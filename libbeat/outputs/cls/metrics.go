package cls

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	writeCLSErrorTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "filebeat",
		Subsystem: "libbeat",
		Name:      "output_http_write_cls_error_total",
		Help:      "The total number of errors when writing output with http to tencent CLS",
	})

	writeCLSRetryTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "filebeat",
		Subsystem: "libbeat",
		Name:      "output_http_write_cls_retry_total",
		Help:      "The total number of retry count when writing output with http to tencent CLS",
	})

	writeCLSLatencyMillis = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "filebeat",
		Subsystem: "libbeat",
		Name:      "output_http_write_cls_latency_millis",
		Help:      "The lantency milliseconds when writing output with http to tencent CLS",
		Buckets:   []float64{10, 20, 30, 40, 50, 100, 150, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 3000},
	})
)