package system

import (
	"fmt"

	"gobot.io/x/gobot"
)

type Optioner interface {
	setDigitalPinToGpiodAccess()
	setSpiToGpioAccess(p gobot.DigitalPinnerProvider, sclkPin, nssPin, mosiPin, misoPin string)
}

func WithDigitalPinGpiodAccess() func(Optioner) {
	return func(s Optioner) {
		s.setDigitalPinToGpiodAccess()
	}
}

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
		sclkPinId:   sclkPin,
		nssPinId:    nssPin,
		mosiPinId:   mosiPin,
		misoPinId:   misoPin,
	}
	gsa := &gpioSpiAccess{cfg: cfg}
	if gsa.isSupported() {
		a.spiAccess = gsa
		if systemDebug {
			fmt.Printf("use gpiod driver for SPI with this config: %s\n", gsa.cfg)
		}
		return
	}
	if systemDebug {
		fmt.Println("gpiod driver not supported for SPI, fallback to periphio")
	}
}
