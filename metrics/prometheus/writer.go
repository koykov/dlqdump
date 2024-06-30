package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricsWriter is a Prometheus implementation of dlqdump.MetricsWriter.
type MetricsWriter struct {
	name string
	prec time.Duration
}

var (
	promSizeIncome, promSizeOutcome, promBytesIncome, promBytesOutcome, promBytesFlush,
	promFail *prometheus.CounterVec

	_ = NewPrometheusMetrics
)

func init() {
	promSizeIncome = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "dlqdump_size_in",
		Help: "Actual queue size.",
	}, []string{"queue"})
	promSizeOutcome = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "dlqdump_size_out",
		Help: "Actual queue size.",
	}, []string{"queue"})

	promBytesIncome = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "dlqdump_bytes_in",
		Help: "How many bytes comes to the queue.",
	}, []string{"queue"})
	promBytesOutcome = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "dlqdump_bytes_out",
		Help: "How many bytes comes to the queue.",
	}, []string{"queue"})
	promBytesFlush = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "dlqdump_bytes_flush",
		Help: "How many bytes flushes from the queue.",
	}, []string{"queue", "reason"})
	promFail = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "dlqdump_fail",
		Help: "Error counters with various reasons.",
	}, []string{"queue", "reason"})

	prometheus.MustRegister(promSizeIncome, promSizeOutcome, promBytesIncome, promBytesOutcome, promBytesFlush, promFail)
}

func NewPrometheusMetrics(name string) *MetricsWriter {
	return NewPrometheusMetricsWP(name, time.Nanosecond)
}

func NewPrometheusMetricsWP(name string, precision time.Duration) *MetricsWriter {
	if precision == 0 {
		precision = time.Nanosecond
	}
	m := &MetricsWriter{
		name: name,
		prec: precision,
	}
	return m
}

func (w MetricsWriter) Dump(size int) {
	promBytesIncome.WithLabelValues(w.name).Add(float64(size))
	promSizeIncome.WithLabelValues(w.name).Inc()
}

func (w MetricsWriter) Flush(reason string, size int) {
	promBytesFlush.WithLabelValues(w.name, reason).Add(float64(size))
}

func (w MetricsWriter) Restore(size int) {
	promBytesOutcome.WithLabelValues(w.name).Add(float64(size))
	promSizeOutcome.WithLabelValues(w.name).Inc()
}

func (w MetricsWriter) Fail(reason string) {
	promFail.WithLabelValues(w.name, reason).Inc()
}
