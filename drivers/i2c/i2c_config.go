package i2c

type i2cConfig struct {
	bus     int
	address int
}

// Config is the interface which describes how a Driver can specify
// optional I2C params such as which I2C bus it wants to use.
type Config interface {
	// WithBus sets which bus to use
	WithBus(bus int)

	// GetBusOrDefault gets which bus to use
	GetBusOrDefault(def int) int

	// WithAddress sets which address to use
	WithAddress(address int)

	// GetAddressOrDefault gets which address to use
	GetAddressOrDefault(def int) int
}

// NewConfig returns a new I2c Config.
func NewConfig() Config {
	return &i2cConfig{bus: BusNotInitialized, address: AddressNotInitialized}
}

// WithBus sets preferred bus to use.
func (i *i2cConfig) WithBus(bus int) {
	i.bus = bus
}

// GetBusOrDefault returns which bus to use, either the one set using WithBus(),
// or the default value which is passed in as the one param.
func (i *i2cConfig) GetBusOrDefault(d int) int {
	if i.bus == BusNotInitialized {
		return d
	}

	return i.bus
}

// WithBus sets which bus to use as a optional param.
func WithBus(bus int) func(Config) {
	return func(i Config) {
		i.WithBus(bus)
	}
}

// WithAddress sets which address to use.
func (i *i2cConfig) WithAddress(address int) {
	i.address = address
}

// GetAddressOrDefault returns which address to use, either
// the one set using WithBus(), or the default value which
// is passed in as the param.
func (i *i2cConfig) GetAddressOrDefault(a int) int {
	if i.address == AddressNotInitialized {
		return a
	}

	return i.address
}

// WithAddress sets which address to use as a optional param.
func WithAddress(address int) func(Config) {
	return func(i Config) {
		i.WithAddress(address)
	}
}
