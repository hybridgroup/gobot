package spi

type spiConfig struct {
	bus   int
	chip  int
	mode  int
	bits  int
	speed int64
}

// NewConfig returns a new SPI Config.
func NewConfig() Config {
	return &spiConfig{
		bus:   NotInitialized,
		chip:  NotInitialized,
		mode:  NotInitialized,
		bits:  NotInitialized,
		speed: NotInitialized,
	}
}

// WithBusNumber sets which bus to use as a optional param.
func WithBusNumber(busNum int) func(Config) {
	return func(s Config) {
		s.SetBusNumber(busNum)
	}
}

// WithChipNumber sets which chip to use as a optional param.
func WithChipNumber(chipNum int) func(Config) {
	return func(s Config) {
		s.SetChipNumber(chipNum)
	}
}

// WithMode sets which mode to use as a optional param.
func WithMode(mode int) func(Config) {
	return func(s Config) {
		s.SetMode(mode)
	}
}

// WithBitCount sets how many bits to use as a optional param.
func WithBitCount(bitCount int) func(Config) {
	return func(s Config) {
		s.SetBitCount(bitCount)
	}
}

// WithSpeed sets what speed to use as a optional param.
func WithSpeed(speed int64) func(Config) {
	return func(s Config) {
		s.SetSpeed(speed)
	}
}

// SetBusNumber sets preferred bus to use.
func (s *spiConfig) SetBusNumber(bus int) {
	s.bus = bus
}

// GetBusNumberOrDefault returns which bus to use, either the one set using WithBus(),
// or the default value which is passed in as the one param.
func (s *spiConfig) GetBusNumberOrDefault(d int) int {
	if s.bus == NotInitialized {
		return d
	}
	return s.bus
}

// SetChipNumber sets preferred chip to use.
func (s *spiConfig) SetChipNumber(chip int) {
	s.chip = chip
}

// GetChipNumberOrDefault returns which chip to use, either the one set using WithChip(),
// or the default value which is passed in as the one param.
func (s *spiConfig) GetChipNumberOrDefault(d int) int {
	if s.chip == NotInitialized {
		return d
	}
	return s.chip
}

// SetMode sets SPI mode to use.
func (s *spiConfig) SetMode(mode int) {
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

// SetBitCount sets how many SPI bits to use.
func (s *spiConfig) SetBitCount(bits int) {
	s.bits = bits
}

// GetBitCountOrDefault returns how many to use, either the one set using WithBits(),
// or the default value which is passed in as the one param.
func (s *spiConfig) GetBitCountOrDefault(d int) int {
	if s.bits == NotInitialized {
		return d
	}
	return s.bits
}

// SetSpeed sets which SPI speed to use.
func (s *spiConfig) SetSpeed(speed int64) {
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
