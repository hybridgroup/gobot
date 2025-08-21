package bleclient

import "time"

// optionApplier needs to be implemented by each configurable option type
type optionApplier interface {
	apply(cfg *configuration)
}

// dropCharacteristicsOnDisconnect is the type for applying drop feature.
type dropCharacteristicsOnDisconnect bool

// debugOption is the type for applying the debug switch on or off.
type debugOption bool

// scanTimeoutOption is the type for applying another timeout than the default 10 min.
type scanTimeoutOption time.Duration

func (o dropCharacteristicsOnDisconnect) String() string {
	return "drop characteristics on disconnect option for BLE client adaptors"
}

func (o debugOption) String() string {
	return "debug option for BLE client adaptors"
}

func (o scanTimeoutOption) String() string {
	return "scan timeout option for BLE client adaptors"
}

func (o dropCharacteristicsOnDisconnect) apply(cfg *configuration) {
	cfg.dropCharacteristicsOnDisconnect = bool(o)
}

func (o debugOption) apply(cfg *configuration) {
	cfg.debug = bool(o)
}

func (o scanTimeoutOption) apply(cfg *configuration) {
	cfg.scanTimeout = time.Duration(o)
}
