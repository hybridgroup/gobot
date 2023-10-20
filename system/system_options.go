package system

import (
	"fmt"

	"gobot.io/x/gobot/v2"
)

// Optioner is the interface for system options. This provides the possibility for change the systems behavior by the
// caller/user when creating the system access, e.g. by "NewAccesser()".
type Optioner interface {
	setDigitalPinToGpiodAccess()
	setSpiToGpioAccess(p gobot.DigitalPinnerProvider, sclkPin, nssPin, mosiPin, misoPin string)
}

// WithDigitalPinGpiodAccess can be used to change the default sysfs implementation for digital pins to the character
// device Kernel ABI. The access is provided by the gpiod package.
func WithDigitalPinGpiodAccess() func(Optioner) {
	return func(s Optioner) {
		s.setDigitalPinToGpiodAccess()
	}
}

// WithSpiGpioAccess can be used to switch the default SPI implementation to GPIO usage.
func WithSpiGpioAccess(p gobot.DigitalPinnerProvider, sclkPin, nssPin, mosiPin, misoPin string) func(Optioner) {
	return func(s Optioner) {
		s.setSpiToGpioAccess(p, sclkPin, nssPin, mosiPin, misoPin)
	}
}

func (a *Accesser) setDigitalPinToGpiodAccess() {
	dpa := &gpiodDigitalPinAccess{fs: a.fs}
	if dpa.isSupported() {
		a.digitalPinAccess = dpa
		if systemDebug {
			fmt.Printf("use gpiod driver for digital pins with this chips: %v\n", dpa.chips)
		}
		return
	}
	if systemDebug {
		fmt.Println("gpiod driver not supported, fallback to sysfs")
	}
}

func (a *Accesser) setSpiToGpioAccess(p gobot.DigitalPinnerProvider, sclkPin, nssPin, mosiPin, misoPin string) {
	cfg := spiGpioConfig{
		pinProvider: p,
		sclkPinID:   sclkPin,
		nssPinID:    nssPin,
		mosiPinID:   mosiPin,
		misoPinID:   misoPin,
	}
	gsa := &gpioSpiAccess{cfg: cfg}
	if gsa.isSupported() {
		a.spiAccess = gsa
		if systemDebug {
			fmt.Printf("use gpio driver for SPI with this config: %s\n", gsa.cfg.String())
		}
		return
	}
	if systemDebug {
		fmt.Println("gpio driver not supported for SPI, fallback to periphio")
	}
}
