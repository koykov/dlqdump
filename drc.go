package dlqdump

// DRC represents dump recycle control.
// DRC may be scheduled (see DRCConfig.CheckInterval).
type DRC struct {
	config *Config
}

func NewDRC(config *Config) (*DRC, error) {
	drc := &DRC{
		config: config.Copy(),
	}
	// todo implement init.
	return drc, nil
}
