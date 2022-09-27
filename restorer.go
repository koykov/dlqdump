package dlqdump

// Restorer represents dump restore handler.
// Restorer may be scheduled (see Config.CheckInterval).
type Restorer struct {
	config *Config
}

func NewRestorer(config *Config) (*Restorer, error) {
	drc := &Restorer{
		config: config.Copy(),
	}
	// todo implement init.
	return drc, nil
}
