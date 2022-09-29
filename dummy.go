package dlqdump

type DummyMetrics struct{}

func (DummyMetrics) Dump(_ string, _ int)     {}
func (DummyMetrics) Flush(_, _ string, _ int) {}
func (DummyMetrics) Restore(_ string, _ int)  {}
