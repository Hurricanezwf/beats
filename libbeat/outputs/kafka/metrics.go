package kafka

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	writeKafkaTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "filebeat",
		Subsystem: "libbeat",
		Name:      "output_kafka_write_total",
		Help:      "The total number  when writing output to kafka",
	})

	writeKafkaErrorTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "filebeat",
		Subsystem: "libbeat",
		Name:      "output_kafka_write_error_total",
		Help:      "The total number of errors when writing output to kafka",
	})

	writeKafkaLatencyMillis = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "filebeat",
		Subsystem: "libbeat",
		Name:      "output_kafka_write_latency_millis",
		Help:      "The lantency milliseconds when writing output to kafka",
		Buckets:   []float64{10, 20, 30, 40, 50, 100, 150, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 3000},
	})
)
