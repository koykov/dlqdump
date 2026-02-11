package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Writer interface {
	Dump(size int)
	Flush(reason string, size int)
	Restore(size int)
	Fail(reason string)
}

// writer is a Prometheus implementation of dlqdump.MetricsWriter.
type writer struct {
	name string
	prec time.Duration
}

var (
	promSizeIncome, promSizeOutcome, promBytesIncome, promBytesOutcome, promBytesFlush,
	promFail *prometheus.CounterVec
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

// NewPrometheusMetrics makes new instance of metrics writer.
// Deprecated: use NewWriter instead.
func NewPrometheusMetrics(name string) Writer {
	return NewWriter(name)
}

// NewPrometheusMetricsWP makes new instance of metrics writer with given precision.
// Deprecated: use NewWriter instead.
func NewPrometheusMetricsWP(name string, precision time.Duration) Writer {
	return NewWriter(name, WithPrecision(precision))
}

// NewWriter makes new instance of metrics writer.
func NewWriter(name string, options ...Option) Writer {
	mw := &writer{name: name}
	for _, fn := range options {
		fn(mw)
	}
	if mw.prec == 0 {
		mw.prec = time.Nanosecond
	}
	return mw
}

func (w *writer) Dump(size int) {
	promBytesIncome.WithLabelValues(w.name).Add(float64(size))
	promSizeIncome.WithLabelValues(w.name).Inc()
}

func (w *writer) Flush(reason string, size int) {
	promBytesFlush.WithLabelValues(w.name, reason).Add(float64(size))
}

func (w *writer) Restore(size int) {
	promBytesOutcome.WithLabelValues(w.name).Add(float64(size))
	promSizeOutcome.WithLabelValues(w.name).Inc()
}

func (w *writer) Fail(reason string) {
	promFail.WithLabelValues(w.name, reason).Inc()
}
