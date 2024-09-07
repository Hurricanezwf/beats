// Note: BY ZWF;
package cls

import (
	"context"
	"fmt"
	"hash/crc64"
	"math"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/outputs/codec"
	"github.com/elastic/beats/v7/libbeat/publisher"
	"github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
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
	log      *logp.Logger
	observer outputs.Observer
	index    string
	codec    codec.Codec
	config   *clsConfig

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

	return &cls{
		log:             logp.NewLogger(logSelector),
		observer:        observer,
		index:           strings.ToLower(index),
		codec:           encoder,
		config:          config,
		nodeIP:          nodeIP,
		producerID:      producerID,
		autoIncrBatchID: 0,
		client:          client,
	}, nil
}

func (c *cls) Close() error {
	return nil
}

func (c *cls) Publish(ctx context.Context, batch publisher.Batch) error {
	events := batch.Events()
	c.observer.NewBatch(len(events))

	packageID := c.generatePackageID()

	var decision string
	var err error
	for i := 0; i < c.config.MaxRetries+1; i++ {
		if decision, err = c.publishEvents(events, packageID); err == nil {
			batch.ACK()
			c.observer.Acked(len(events))
			return nil
		}
		time.Sleep(time.Second)
	}

	switch decision {
	case "ack":
		batch.ACK()
		c.observer.Acked(len(events))
	case "drop":
		batch.Drop()
		c.observer.Dropped(len(events))
	case "retry":
		batch.Retry()
		c.observer.Failed(len(events))
	}
	return fmt.Errorf("unknown decision %s returned", decision)
}

func (c *cls) publishEvents(events []publisher.Event, packageID string) (decision string, err error) {
	_ = c.logGroupListFrom(events)

	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()

	//c.client.Send(ctx, c.config.Topic)

	return "ack", nil
}

func (c *cls) String() string {
	return "cls"
}

// event 是 flattern 过的格式:
// {"agent.ephemeral_id":"d49c7cb0-f9ae-4696-b44d-4bb68e7ce296","agent.id":"1a300adf-207e-4ef5-ab29-4dc4d875f655","agent.name":"zwf","agent.type":"filebeat","agent.version":"8.14.4","ecs.version":"8.0.0","host.name":"zwf","input.type":"log","log.file.path":"/home/zwf/workspace/experiment/beats/filebeat/y.log","log.offset":88,"message":"hello","metadata.version":"8.14.4","timestamp":"2024-09-07T17:00:11.281090418+08:00"}
func (c *cls) logGroupListFrom(events []publisher.Event) LogGroupList {
	//fileGrp := make(map[string]LogGroup)
	for _, e := range events {
		fields := e.Content.Fields.StringToPrint()
		fmt.Printf("fields=%s\n", fields)
		meta := e.Content.Meta.StringToPrint()
		fmt.Printf("meta=%s\n", meta)
		ts := e.Content.Timestamp
		fmt.Printf("ts=%s\n", ts.String())
		return LogGroupList{}
	}
	return LogGroupList{}
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
