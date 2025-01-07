package system

import (
	"gobot.io/x/gobot/v2"
)

// accesserOptionApplier is the interface for system options. This provides the possibility for change the systems
// behavior by the caller/user when creating the system access, e.g. by "NewAccesser().Add..".
// The interface needs to be implemented by each configurable option type.
type AccesserOptionApplier interface {
	apply(cfg *accesserConfiguration)
}

type systemAccesserDebugOption bool

type systemDigitalPinDebugOption bool

type systemSpiDebugOption bool

type systemUseDigitalPinSysfsOption bool

type systemUseSpiGpioOption spiGpioConfig

// WithSystemAccesserDebug can be used to switch on debug messages.
func WithSystemAccesserDebug() systemAccesserDebugOption {
	return systemAccesserDebugOption(true)
}

// WithDigitalPinDebug can be used to switch on debug messages for digital pins.
func WithDigitalPinDebug() systemDigitalPinDebugOption {
	return systemDigitalPinDebugOption(true)
}

// WithSpiDebug can be used to switch on debug messages for SPI.
func WithSpiDebug() systemSpiDebugOption {
	return systemSpiDebugOption(true)
}

// WithDigitalPinSysfsAccess can be used to change the default character device implementation for digital pins to the
// legacy sysfs Kernel ABI.
func WithDigitalPinSysfsAccess() systemUseDigitalPinSysfsOption {
	return systemUseDigitalPinSysfsOption(true)
}

// WithDigitalPinCdevAccess can be used to change the default sysfs implementation for digital pins in old platforms to
// test the character device Kernel ABI. The access is provided by the go-gpiocdev package.
func WithDigitalPinCdevAccess() systemUseDigitalPinSysfsOption {
	return systemUseDigitalPinSysfsOption(false)
}

// WithSpiGpioAccess can be used to switch the default SPI implementation to GPIO usage.
func WithSpiGpioAccess(p gobot.DigitalPinnerProvider, sclkPin, ncsPin, sdoPin, sdiPin string) systemUseSpiGpioOption {
	o := systemUseSpiGpioOption{
		pinProvider: p,
		sclkPinID:   sclkPin,
		ncsPinID:    ncsPin,
		sdoPinID:    sdoPin,
		sdiPinID:    sdiPin,
	}

	return o
}

func (o systemAccesserDebugOption) String() string {
	return "switch on system accesser debugging option"
}

func (o systemDigitalPinDebugOption) String() string {
	return "switch on system digital pin debugging option"
}

func (o systemSpiDebugOption) String() string {
	return "switch on system SPI debugging option"
}

func (o systemUseDigitalPinSysfsOption) String() string {
	return "system accesser use sysfs vs. cdev for digital pins option"
}

func (o systemUseSpiGpioOption) String() string {
	return "system accesser use discrete GPIOs for SPI option"
}

func (o systemAccesserDebugOption) apply(cfg *accesserConfiguration) {
	cfg.debug = bool(o)
}

func (o systemDigitalPinDebugOption) apply(cfg *accesserConfiguration) {
	cfg.debugDigitalPin = bool(o)
}

func (o systemSpiDebugOption) apply(cfg *accesserConfiguration) {
	cfg.debugSpi = bool(o)
}

func (o systemUseDigitalPinSysfsOption) apply(cfg *accesserConfiguration) {
	c := bool(o)
	cfg.useGpioSysfs = &c
}

func (o systemUseSpiGpioOption) apply(cfg *accesserConfiguration) {
	c := spiGpioConfig(o)
	cfg.spiGpioConfig = &c
}
