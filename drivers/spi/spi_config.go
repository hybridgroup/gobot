package spi

type spiConfig struct {
	bus   int
	chip  int
	mode  int
	bits  int
	speed int64
}

// Config is the interface which describes how a Driver can specify
// optional SPI params such as which SPI bus it wants to use.
type Config interface {
	// WithBus sets which bus to use
	WithBus(bus int)

	// GetBusOrDefault gets which bus to use
	GetBusOrDefault(def int) int

	// WithChip sets which chip to use
	WithChip(chip int)

	// GetChipOrDefault gets which chip to use
	GetChipOrDefault(def int) int

	// WithMode sets which mode to use
	WithMode(mode int)

	// GetModeOrDefault gets which mode to use
	GetModeOrDefault(def int) int

	// WithBIts sets how many bits to use
	WithBits(bits int)

	// GetBitsOrDefault gets how many bits to use
	GetBitsOrDefault(def int) int

	// WithSpeed sets which speed to use (in Hz)
	WithSpeed(speed int64)

	// GetSpeedOrDefault gets which speed to use (in Hz)
	GetSpeedOrDefault(def int64) int64
}

// NewConfig returns a new SPI Config.
func NewConfig() Config {
	return &spiConfig{
		bus:   NotInitialized,
		chip:  NotInitialized,
		mode:  NotInitialized,
		bits:  NotInitialized,
		speed: NotInitialized}
}

// WithBus sets preferred bus to use.
func (s *spiConfig) WithBus(bus int) {
	s.bus = bus
}

// GetBusOrDefault returns which bus to use, either the one set using WithBus(),
// or the default value which is passed in as the one param.
func (s *spiConfig) GetBusOrDefault(d int) int {
	if s.bus == NotInitialized {
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

// WithChip sets preferred chip to use.
func (s *spiConfig) WithChip(chip int) {
	s.chip = chip
}

// GetChipOrDefault returns which chip to use, either the one set using WithChip(),
// or the default value which is passed in as the one param.
func (s *spiConfig) GetChipOrDefault(d int) int {
	if s.chip == NotInitialized {
		return d
	}
	return s.chip
}

// WithChip sets which chip to use as a optional param.
func WithChip(chip int) func(Config) {
	return func(s Config) {
		s.WithChip(chip)
	}
}

// WithMode sets SPI mode to use.
func (s *spiConfig) WithMode(mode int) {
	s.mode = mode
}

// GetModeOrDefault returns which mode to use, either the one set using WithChip(),
// or the default value which is passed in as the one param.
func (s *spiConfig) GetModeOrDefault(d int) int {
	if s.mode == NotInitialized {
		return d
	}
	return s.mode
}

// WithMode sets which mode to use as a optional param.
func WithMode(mode int) func(Config) {
	return func(s Config) {
		s.WithMode(mode)
	}
}

// WithBits sets how many SPI bits to use.
func (s *spiConfig) WithBits(bits int) {
	s.bits = bits
}

// GetBitsOrDefault returns how many to use, either the one set using WithBits(),
// or the default value which is passed in as the one param.
func (s *spiConfig) GetBitsOrDefault(d int) int {
	if s.bits == NotInitialized {
		return d
	}
	return s.bits
}

// WithBits sets how many bits to use as a optional param.
func WithBits(bits int) func(Config) {
	return func(s Config) {
		s.WithBits(bits)
	}
}

// WithSpeed sets which SPI speed to use.
func (s *spiConfig) WithSpeed(speed int64) {
	s.speed = speed
}

// GetSpeedOrDefault returns what speed to use, either the one set using WithSpeed(),
// or the default value which is passed in as the one param.
func (s *spiConfig) GetSpeedOrDefault(d int64) int64 {
	if s.speed == NotInitialized {
		return d
	}
	return s.speed
}

// WithSpeed sets what speed to use as a optional param.
func WithSpeed(speed int64) func(Config) {
	return func(s Config) {
		s.WithSpeed(speed)
	}
}
