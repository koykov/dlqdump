package dlqdump

type DummyMetrics struct{}

func (DummyMetrics) DumpPut(_ string, _ int)      {}
func (DummyMetrics) DumpFlush(_, _ string, _ int) {}
func (DummyMetrics) DumpRestore(_ string, _ int)  {}
