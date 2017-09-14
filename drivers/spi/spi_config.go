package spi

type spiConfig struct {
	bus int
}

// Config is the interface which describes how a Driver can specify
// optional SPI params such as which SPI bus it wants to use.
type Config interface {
	// WithBus sets which bus to use
	WithBus(bus int)

	// GetBusOrDefault gets which bus to use
	GetBusOrDefault(def int) int
}

// NewConfig returns a new SPI Config.
func NewConfig() Config {
	return &spiConfig{bus: BusNotInitialized}
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
