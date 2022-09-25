package dlqdump

// DRC represents [scheduled] dump recycle control.
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
