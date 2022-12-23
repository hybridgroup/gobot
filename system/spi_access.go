package system

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
)

type periphioSpiAccess struct{}

type gpiodSpiAccess struct {
	fs   filesystem
	sclk int
	ssz  int
	mosi int
	miso int
}

func (*periphioSpiAccess) createDevice(busNum, chipNum, mode, bits int, maxSpeed int64) (gobot.SpiSystemDevicer, error) {
	return newSpiPeriphIo(busNum, chipNum, mode, bits, maxSpeed)
}

func (*periphioSpiAccess) isSupported() bool {
	/*
		devices, err := h.fs.find("/dev", "spidev")
		if err != nil || len(devices) == 0 {
			return false
		}
	*/
	return true
}

func (gsa *gpiodSpiAccess) createDevice(busNum, chipNum, mode, bits int, maxSpeed int64) (gobot.SpiSystemDevicer, error) {
	// maxSpeed is given in Hz, tclk is half the cycle time, tclk=1/(2*f)
	// default is 1 MHz => tclk = 500 nano seconds, range can be reach 10MHz
	// tclk[ns]=1 000 000 000/(2*maxSpeed)
	tclk := time.Duration(1000000000/2/maxSpeed) * time.Nanosecond
	chipName := fmt.Sprintf("gpiochip%d", busNum)
	return newSpiGpiod(chipName, gsa.sclk, gsa.ssz, gsa.mosi, gsa.miso, tclk)
}

func (gsa *gpiodSpiAccess) isSupported() bool {
	chips, err := gsa.fs.find("/dev", "gpiochip")
	if err != nil || len(chips) == 0 {
		return false
	}
	return true
}
