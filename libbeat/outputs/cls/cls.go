// Note: BY ZWF;
package cls

import (
	"context"
	"errors"
	"fmt"
	"hash/crc64"
	"math"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/outputs/codec"
	"github.com/elastic/beats/v7/libbeat/publisher"
	"github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"

	"github.com/tidwall/gjson"
)

const (
	logSelector = "cls"
)

func init() {
	outputs.RegisterType("cls", makeCLS)
}

func makeCLS(
	_ outputs.IndexManager,
	beat beat.Info,
	observer outputs.Observer,
	cfg *config.C,
) (outputs.Group, error) {
	log := logp.NewLogger(logSelector)
	log.Debug("initialize cls output")

	config, err := readConfig(cfg)
	if err != nil {
		return outputs.Fail(fmt.Errorf("failed to read output.cls config, %w", err))
	}

	codec, err := codec.CreateEncoder(beat, config.Codec)
	if err != nil {
		return outputs.Fail(err)
	}

	client, err := newCLSClient(observer, beat.IndexPrefix, codec, config)
	if err != nil {
		return outputs.Fail(fmt.Errorf("failed to new cls client, %w", err))
	}

	retry := 0
	if config.MaxRetries < 0 {
		retry = -1
	}

	return outputs.Success(config.Queue, config.BulkMaxSize, retry, nil, client)
}

type cls struct {
	ctx      context.Context
	cancel   context.CancelFunc
	log      *logp.Logger
	observer outputs.Observer
	index    string
	codec    codec.Codec
	config   *clsConfig
	queue    chan publisher.Batch
	wg       *sync.WaitGroup

	nodeIP          string
	producerID      string
	autoIncrBatchID uint64
	client          *CLSHTTPClient
}

func newCLSClient(observer outputs.Observer, index string, encoder codec.Codec, config *clsConfig) (*cls, error) {
	// 计算 producer ID, 为了后续支持在腾讯云上进行日志上下文查询.
	nodeIP := os.Getenv(config.NodeIPEnv)
	if nodeIP == "" {
		return nil, fmt.Errorf("empty nodeIP from env %s", config.NodeIPEnv)
	}
	producerID := generateProducerHash(fmt.Sprintf("%s-%d", nodeIP, time.Now().UnixNano()))

	// new cls http client;
	client, err := NewCLSHTTPClient(config.Endpoint, config.AccessKey, config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to new cls http client, %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &cls{
		ctx:             ctx,
		cancel:          cancel,
		log:             logp.NewLogger(logSelector),
		observer:        observer,
		index:           strings.ToLower(index),
		codec:           encoder,
		config:          config,
		queue:           make(chan publisher.Batch, config.WriteQueueSize),
		wg:              &sync.WaitGroup{},
		nodeIP:          nodeIP,
		producerID:      producerID,
		autoIncrBatchID: 0,
		client:          client,
	}

	c.startWorkers()

	return c, nil
}

func (c *cls) startWorkers() {
	c.wg.Add(c.config.WriteConcurrency)
	for i := 0; i < c.config.WriteConcurrency; i++ {
		go func() {
			defer c.wg.Done()
			for batch := range c.queue {
				_ = c.publish(batch)
			}
		}()
	}
}

func (c *cls) Close() error {
	c.log.Info("output.cls is closing")
	// 1.阻断publish;
	if c.cancel != nil {
		c.cancel()
	}
	time.Sleep(time.Second)
	// 2. 关闭队列;
	if c.queue != nil {
		close(c.queue)
	}
	// 3. 等待所有worker安全退出;
	c.wg.Wait()
	// 4. 关闭成功;
	c.log.Info("output.cls closed")
	return nil
}

func (c *cls) Publish(ctx context.Context, batch publisher.Batch) error {
	select {
	case <-c.ctx.Done():
		batch.Cancelled()
		return errors.New("output.cls was closed")
	default:
	}

	select {
	case c.queue <- batch:
		return nil
	case <-ctx.Done():
		// 如果队列满了, 卡到超时, 设置状态为取消, 不计次重试;
		batch.Cancelled()
		return errors.New("queue is full")
	case <-c.ctx.Done():
		batch.Cancelled()
		return errors.New("output.cls was closed")
	}
}

func (c *cls) publish(batch publisher.Batch) error {
	events := batch.Events()

	// 将所有事件编码成 flatten 模式;
	var flatternJSONEvents []gjson.Result
	for _, e := range events {
		b, err := c.codec.Encode(c.index, &e.Content)
		if err != nil {
			c.log.Warnf("failed to encode event content, %v", err)
			continue
		}
		flatternJSONEvents = append(flatternJSONEvents, gjson.ParseBytes(b))
	}

	c.observer.NewBatch(len(events))

	packageID := c.generatePackageID()

	var decision string
	var fastRecover bool
	var err error
	var start = time.Now()
	for i := 0; i < c.config.MaxRetries+1; i++ {
		if i > 0 {
			writeCLSRetryTotal.Inc()
		}
		decision, fastRecover, err = c.publishEvents(flatternJSONEvents, packageID)
		if err == nil || !fastRecover {
			break
		}
		time.Sleep(time.Second)
	}

	switch decision {
	case "ack":
		batch.ACK()
		c.observer.Acked(len(events))
		c.observer.ReportLatency(time.Since(start))
		writeCLSLatencyMillis.Observe(float64(time.Since(start).Milliseconds()))
		return err
	case "drop":
		batch.Drop()
		c.observer.Dropped(len(events))
		return err
	case "retry":
		batch.Retry()
		c.observer.Failed(len(events))
		writeCLSRetryTotal.Inc()
		return err
	}
	return fmt.Errorf("unknown decision %s returned", decision)
}

func (c *cls) publishEvents(flatternJSONEvents []gjson.Result, packageID string) (decision string, fastRecover bool, err error) {
	logGroupList := c.logGroupListFrom(flatternJSONEvents, packageID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clserr := c.client.Send(ctx, c.config.Topic, logGroupList)
	if clserr == nil {
		return "ack", false, nil
	}
	writeCLSErrorTotal.Inc()

	if clserr.HTTPCode < 0 {
		if clserr.Code == TEMPORARY_ERROR {
			return "retry", true, errors.New(clserr.Message)
		}
		return "drop", false, fmt.Errorf("error code: %s, message: %s", clserr.Code, clserr.Message)
	}

	nxx := clserr.HTTPCode / 100
	switch nxx {
	case 2:
		return "ack", false, nil
	case 3:
		// redirect;
		return "retry", false, fmt.Errorf("cls server return status code %d, errorCode: %s, msg: %s, requestID: %s", clserr.HTTPCode, clserr.Code, clserr.Message, clserr.RequestID)
	case 4:
		// client error;
		return "retry", false, fmt.Errorf("cls server return status code %d, errorCode: %s, msg: %s, requestID: %s", clserr.HTTPCode, clserr.Code, clserr.Message, clserr.RequestID)
	case 5:
		// server error;
		return "retry", true, fmt.Errorf("cls server return status code %d, errorCode: %s, msg: %s, requestID: %s", clserr.HTTPCode, clserr.Code, clserr.Message, clserr.RequestID)
	}
	return "drop", false, fmt.Errorf("invalid http status code %d returned from cls, errorCode: %s, msg: %s, requestID: %s", clserr.HTTPCode, clserr.Code, clserr.Message, clserr.RequestID)
}

func (c *cls) String() string {
	return "cls"
}

// event 是 flattern 过的格式:
// {"agent.ephemeral_id":"d49c7cb0-f9ae-4696-b44d-4bb68e7ce296","agent.id":"1a300adf-207e-4ef5-ab29-4dc4d875f655","agent.name":"zwf","agent.type":"filebeat","agent.version":"8.14.4","ecs.version":"8.0.0","host.name":"zwf","input.type":"log","log.file.path":"/home/zwf/workspace/experiment/beats/filebeat/y.log","log.offset":88,"message":"hello","metadata.version":"8.14.4","timestamp":1762282382349234324}
func (c *cls) logGroupListFrom(flatternJSONEvents []gjson.Result, packageID string) LogGroupList {
	nowTS := time.Now().UnixNano()
	fileGrp := make(map[string]*LogGroup)
	for _, e := range flatternJSONEvents {
		filepath := e.Get("log\\.file\\.path").String()
		if filepath == "" {
			c.log.Errorf("could not find log.file.path from input flatten JSON event, json str: %s", e.String())
			continue
		}
		m := e.Map()
		contents := make([]*Log_Content, 0, len(m))
		for k, v := range m {
			contents = append(contents, &Log_Content{
				Key:   stringPtr(k),
				Value: stringPtr(v.String()),
			})
		}

		grp := fileGrp[filepath]
		if grp == nil {
			grp = &LogGroup{}
		}
		grp.Logs = append(grp.Logs, &Log{
			Time:        int64Ptr(e.Get("timestamp").Int()),
			Contents:    contents,
			CollectTime: int64Ptr(nowTS),
		})
		if grp.ContextFlow == nil {
			grp.ContextFlow = stringPtr(packageID)
		}
		if grp.Filename == nil {
			grp.Filename = stringPtr(filepath)
		}
		if grp.Source == nil {
			grp.Source = stringPtr(c.nodeIP)
		}
		if grp.Hostname == nil {
			grp.Hostname = stringPtr(c.nodeIP)
		}
		fileGrp[filepath] = grp
	}

	var logGroupList LogGroupList
	for _, grp := range fileGrp {
		logGroupList.LogGroupList = append(logGroupList.LogGroupList, grp)
	}
	return logGroupList
}

func (c *cls) generatePackageID() string {
	batchID := atomic.AddUint64(&c.autoIncrBatchID, 1)
	if batchID >= math.MaxUint64-1 {
		batchID = 0
	}
	return strings.ToUpper(fmt.Sprintf("%s-%x", c.producerID, batchID))
}

func generateProducerHash(str string) string {
	table := crc64.MakeTable(crc64.ECMA)
	hash := crc64.Checksum([]byte(str), table)
	hashString := fmt.Sprintf("%016x", hash)
	return strings.ToUpper(hashString)
}

func stringPtr(str string) *string {
	strcopy := str
	return &strcopy
}

func int64Ptr(v int64) *int64 {
	vcopy := v
	return &vcopy
}
