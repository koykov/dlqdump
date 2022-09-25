package dlqdump

// DRC represents dump recycle control.
// DRC may be scheduled (see DRConfig.CheckInterval).
type DRC struct {
	config *DRConfig
}

func NewDRC(config *DRConfig) (*DRC, error) {
	drc := &DRC{
		config: config.Copy(),
	}
	// todo implement init.
	return drc, nil
}
