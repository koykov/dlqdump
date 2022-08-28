package dlqdump

type DummyMetrics struct{}

func (DummyMetrics) QueuePut(_ string, _ int)      {}
func (DummyMetrics) QueueFlush(_, _ string, _ int) {}
