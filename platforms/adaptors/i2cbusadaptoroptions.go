package adaptors

// I2CBusOptionApplier is the interface for I2C bus adaptors options. This provides the possibility for change the
// platform behavior by the user when creating the platform, e.g. by "NewAdaptor()".
// The interface needs to be implemented by each configurable option type.
type I2CBusOptionApplier interface {
	apply(cfg *i2cBusConfiguration)
}

// i2cBusDebugOption is the type to switch on I2C related debug messages.
type i2cBusDebugOption bool

func (o i2cBusDebugOption) String() string {
	return "switch on debugging for I2C option"
}

func (o i2cBusDebugOption) apply(cfg *i2cBusConfiguration) {
	cfg.debug = bool(o)
}
