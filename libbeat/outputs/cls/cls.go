package cls

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
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

	httpcli *http.Client
	//writer *bufio.Writer
	//out    *os.File
}

func newCLSClient(observer outputs.Observer, index string, encoder codec.Codec, config *clsConfig) (*cls, error) {
	return &cls{
		log:      logp.NewLogger(logSelector),
		observer: observer,
		index:    strings.ToLower(index),
		codec:    encoder,
		config:   config,
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

	dropped := 0
	for i := range events {
		ok := c.publishEvent(&events[i])
		if !ok {
			dropped++
		}
	}

	//c.writer.Flush()
	batch.ACK()

	c.observer.Dropped(dropped)
	c.observer.Acked(len(events) - dropped)

	return nil
}

func (c *cls) publishEvent(event *publisher.Event) bool {
	serializedEvent, err := c.codec.Encode(c.index, &event.Content)
	if err != nil {
		if !event.Guaranteed() {
			return false
		}

		c.log.Errorf("Unable to encode event: %+v", err)
		c.log.Debugf("Failed event: %v", event)
		return false
	}

	os.Stderr.WriteString("\n-----------------------------------------------\n")
	if _, err = os.Stderr.Write(serializedEvent); err != nil {
		c.observer.WriteError(err)
		c.log.Errorf("Error when appending newline to event: %+v", err)
		return false
	}

	c.observer.WriteBytes(len(serializedEvent) + 1)
	return true
}

func (c *cls) String() string {
	return "cls"
}
