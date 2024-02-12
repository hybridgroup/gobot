package serialport

// optionApplier needs to be implemented by each configurable option type
type optionApplier interface {
	apply(cfg *configuration)
}

// nameOption is the type for applying another name to the configuration
type nameOption string

// baudRateOption is the type for applying another baud rate than the default 115200
type baudRateOption int

func (o nameOption) String() string {
	return "name option for Serial Port adaptors"
}

func (o baudRateOption) String() string {
	return "baud rate option for Serial Port adaptors"
}

func (o nameOption) apply(cfg *configuration) {
	cfg.name = string(o)
}

func (o baudRateOption) apply(cfg *configuration) {
	cfg.baudRate = int(o)
}
