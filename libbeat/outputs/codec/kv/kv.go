package kv

import (
	"encoding/json"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/outputs/codec"
	"github.com/elastic/elastic-agent-libs/config"
)

func init() {
	codec.RegisterType("kv", func(info beat.Info, cfg *config.C) (codec.Codec, error) {
		return New(info.Version), nil
	})
}

// Encoder for serializing a beat.Event to json.
type Encoder struct {
	version string
}

// New creates a new json Encoder.
func New(version string) *Encoder {
	return &Encoder{
		version: version,
	}
}

// Encode serializes a beat event to JSON. It adds additional metadata in the
// {"@timestamp":"2024-09-07T08:02:54.662Z","@metadata":{"beat":"filebeat","type":"_doc","version":"8.14.4"},"ecs":{"version":"8.0.0"},"log":{"file":{"path":"/home/zwf/workspace/experiment/beats/filebeat/x.log"},"offset":0},"message":"hello world","input":{"type":"log"},"host":{"name":"zwf"},"agent":{"id":"1a300adf-207e-4ef5-ab29-4dc4d875f655","name":"zwf","type":"filebeat","version":"8.14.4","ephemeral_id":"270c28c5-13e7-43f3-b27b-22382d331408"}}
func (e *Encoder) Encode(index string, event *beat.Event) ([]byte, error) {
	kvs := make(map[string]any)
	kvs["timestamp"] = event.Timestamp.In(time.FixedZone("UTC+8", 8*3600)).Format(time.RFC3339Nano)
	kvs["metadata.version"] = e.version
	for k, v := range event.Fields.Flatten() {
		kvs[k] = v
	}
	return json.Marshal(kvs)
}
