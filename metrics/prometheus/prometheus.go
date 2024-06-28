package dlqdump

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMetrics is a Prometheus implementation of dlqdump.MetricsWriter.
type PrometheusMetrics struct {
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

func NewPrometheusMetrics(name string) *PrometheusMetrics {
	return NewPrometheusMetricsWP(name, time.Nanosecond)
}

func NewPrometheusMetricsWP(name string, precision time.Duration) *PrometheusMetrics {
	if precision == 0 {
		precision = time.Nanosecond
	}
	m := &PrometheusMetrics{
		name: name,
		prec: precision,
	}
	return m
}

func (m PrometheusMetrics) Dump(size int) {
	promBytesIncome.WithLabelValues(m.name).Add(float64(size))
	promSizeIncome.WithLabelValues(m.name).Inc()
}

func (m PrometheusMetrics) Flush(reason string, size int) {
	promBytesFlush.WithLabelValues(m.name, reason).Add(float64(size))
}

func (m PrometheusMetrics) Restore(size int) {
	promBytesOutcome.WithLabelValues(m.name).Add(float64(size))
	promSizeOutcome.WithLabelValues(m.name).Inc()
}

func (m PrometheusMetrics) Fail(reason string) {
	promFail.WithLabelValues(m.name, reason).Inc()
}
