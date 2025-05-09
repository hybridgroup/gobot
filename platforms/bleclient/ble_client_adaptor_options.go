package bleclient

import "time"

// optionApplier needs to be implemented by each configurable option type
type optionApplier interface {
	apply(cfg *configuration)
}

// debugOption is the type for applying the debug switch on or off.
type debugOption bool

// scanTimeoutOption is the type for applying another timeout than the default 10 min.
type scanTimeoutOption time.Duration

func (o debugOption) String() string {
	return "debug option for BLE client adaptors"
}

func (o scanTimeoutOption) String() string {
	return "scan timeout option for BLE client adaptors"
}

func (o debugOption) apply(cfg *configuration) {
	cfg.debug = bool(o)
}

func (o scanTimeoutOption) apply(cfg *configuration) {
	cfg.scanTimeout = time.Duration(o)
}
