package spi

type spiConfig struct {
	bus     int
	address int
}

// Config is the interface which describes how a Driver can specify
// optional SPI params such as which SPI bus it wants to use.
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

// NewConfig returns a new SPI Config.
func NewConfig() Config {
	return &spiConfig{bus: BusNotInitialized, address: AddressNotInitialized}
}

// WithBus sets preferred bus to use.
func (s *spiConfig) WithBus(bus int) {
	s.bus = bus
}

// GetBusOrDefault returns which bus to use, either the one set using WithBus(),
// or the default value which is passed in as the one param.
func (s *spiConfig) GetBusOrDefault(d int) int {
	if s.bus == BusNotInitialized {
		return d
	}
	return s.bus
}

// WithBus sets which bus to use as a optional param.
func WithBus(bus int) func(Config) {
	return func(s Config) {
		s.WithBus(bus)
	}
}

// WithAddress sets which address to use.
func (s *spiConfig) WithAddress(address int) {
	s.address = address
}

// GetAddressOrDefault returns which address to use, either
// the one set using WithBus(), or the default value which
// is passed in as the param.
func (s *spiConfig) GetAddressOrDefault(a int) int {
	if s.address == AddressNotInitialized {
		return a
	}
	return s.address
}

// WithAddress sets which address to use as a optional param.
func WithAddress(address int) func(Config) {
	return func(s Config) {
		s.WithAddress(address)
	}
}
