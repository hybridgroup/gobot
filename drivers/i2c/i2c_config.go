package i2c

type i2cConfig struct {
	bus     int
	address int
}

// NewConfig returns a new I2c Config.
func NewConfig() Config {
	return &i2cConfig{bus: BusNotInitialized, address: AddressNotInitialized}
}

// WithBus sets which bus to use as a optional param.
func WithBus(bus int) func(Config) {
	return func(i Config) {
		i.SetBus(bus)
	}
}

// WithAddress sets which address to use as a optional param.
func WithAddress(address int) func(Config) {
	return func(i Config) {
		i.SetAddress(address)
	}
}

// SetBus sets preferred bus to use.
func (i *i2cConfig) SetBus(bus int) {
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

// SetAddress sets which address to use.
func (i *i2cConfig) SetAddress(address int) {
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
