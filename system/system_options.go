package system

import "fmt"

type systemOptioner interface {
	setDigitalPinAccess(string)
	setSpiGpiodAccess(sclk, ssz, mosi, miso int)
}

func WithDigitalPinAccess(val string) func(systemOptioner) {
	return func(s systemOptioner) {
		s.setDigitalPinAccess(val)
	}
}

func WithSpiGiodAccess(sclk, ssz, mosi, miso int) func(systemOptioner) {
	return func(s systemOptioner) {
		s.setSpiGpiodAccess(sclk, ssz, mosi, miso)
	}
}

func (a *Accesser) setDigitalPinAccess(digitalPinAccess string) {
	if digitalPinAccess == "" {
		digitalPinAccess = "sysfs"
	}
	if digitalPinAccess != "sysfs" {
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
	a.digitalPinAccess = &sysfsDigitalPinAccess{fs: a.fs}
}

func (a *Accesser) setSpiGpiodAccess(sclk, ssz, mosi, miso int) {
	gsa := &gpiodSpiAccess{
		sclk: sclk,
		ssz:  ssz,
		mosi: mosi,
		miso: miso,
		fs:   a.fs,
	}
	if gsa.isSupported() {
		a.spiAccess = gsa
		if systemDebug {
			fmt.Printf("use gpiod driver for SPI with this pins:\n")
			fmt.Printf("sclk: %d, ssz: %d, mosi: %d, miso: %d\n", gsa.sclk, gsa.ssz, gsa.mosi, gsa.miso)
		}
		return
	}
	if systemDebug {
		fmt.Println("gpiod driver not supported for SPI, fallback to periphio")
	}
}
