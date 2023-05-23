package system

import (
	"gobot.io/x/gobot/v2"
)

type periphioSpiAccess struct {
	fs filesystem
}

type gpioSpiAccess struct {
	cfg spiGpioConfig
}

func (*periphioSpiAccess) createDevice(busNum, chipNum, mode, bits int, maxSpeed int64) (gobot.SpiSystemDevicer, error) {
	return newSpiPeriphIo(busNum, chipNum, mode, bits, maxSpeed)
}

func (psa *periphioSpiAccess) isSupported() bool {
	devices, err := psa.fs.find("/dev", "spidev")
	if err != nil || len(devices) == 0 {
		return false
	}
	return true
}

func (gsa *gpioSpiAccess) createDevice(busNum, chipNum, mode, bits int, maxSpeed int64) (gobot.SpiSystemDevicer, error) {
	return newSpiGpio(gsa.cfg, maxSpeed)
}

func (gsa *gpioSpiAccess) isSupported() bool {
	return true
}
