package adaptors

import "gobot.io/x/gobot/v2/system"

// SpiBusOptionApplier is the interface for spi bus adaptor options. This provides the possibility for change the
// platform behavior by the user when creating the platform, e.g. by "NewAdaptor()".
// The interface needs to be implemented by each configurable option type.
type SpiBusOptionApplier interface {
	apply(cfg *spiBusConfiguration)
}

// spiBusDebugOption is the type to switch on SPI related debug messages.
type spiBusDebugOption bool

// spiBusDigitalPinsForSystemSpiOption is the type to switch the default SPI implementation to GPIO usage
type spiBusDigitalPinsForSystemSpiOption struct {
	sclkPin string
	ncsPin  string
	sdoPin  string
	sdiPin  string
}

func (o spiBusDebugOption) String() string {
	return "switch on debugging for SPI option"
}

func (o spiBusDigitalPinsForSystemSpiOption) String() string {
	return "use digital pins for SPI option"
}

func (o spiBusDebugOption) apply(cfg *spiBusConfiguration) {
	cfg.debug = bool(o)
	cfg.systemOptions = append(cfg.systemOptions, system.WithSpiDebug())
}

func (o spiBusDigitalPinsForSystemSpiOption) apply(cfg *spiBusConfiguration) {
	cfg.systemOptions = append(cfg.systemOptions, system.WithSpiGpioAccess(cfg.spiGpioPinnerProvider, o.sclkPin, o.ncsPin,
		o.sdoPin, o.sdiPin))
}
