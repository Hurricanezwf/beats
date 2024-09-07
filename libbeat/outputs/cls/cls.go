// Note: BY ZWF;
package cls

import (
	"context"
	"fmt"
	"hash/crc64"
	"math"
	"net"
	"net/http"
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

	producerID      string
	autoIncrBatchID uint64
	httpcli         *http.Client
}

func newCLSClient(observer outputs.Observer, index string, encoder codec.Codec, config *clsConfig) (*cls, error) {
	// 计算 producer ID, 为了后续支持在腾讯云上进行日志上下文查询.
	nodeIP := os.Getenv(config.NodeIPEnv)
	if nodeIP == "" {
		return nil, fmt.Errorf("empty nodeIP from env %s", config.NodeIPEnv)
	}
	producerID := generateProducerHash(fmt.Sprintf("%s-%d", nodeIP, time.Now().UnixNano()))

	return &cls{
		log:             logp.NewLogger(logSelector),
		observer:        observer,
		index:           strings.ToLower(index),
		codec:           encoder,
		config:          config,
		producerID:      producerID,
		autoIncrBatchID: 0,
		httpcli: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
					dialer := &net.Dialer{
						Timeout:   5 * time.Second,
						KeepAlive: 30 * time.Second,
					}
					return dialer.DialContext(ctx, network, addr)
				},
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}, nil
}

func (c *cls) Close() error {
	// TODO: 按需 close;
	return nil
}

func (c *cls) Publish(_ context.Context, batch publisher.Batch) error {
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
	_ = events
	_ = packageID
	// TODO:
	return "ack", nil
}

func (c *cls) String() string {
	return "cls"
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
