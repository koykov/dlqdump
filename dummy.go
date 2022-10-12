package dlqdump

type DummyMetrics struct{}

func (DummyMetrics) Dump(_ int)            {}
func (DummyMetrics) Flush(_ string, _ int) {}
func (DummyMetrics) Restore(_ int)         {}
func (DummyMetrics) Fail(_ string)         {}
