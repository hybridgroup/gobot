package system

import "gobot.io/x/gobot"

type periphioSpiAccess struct{}

func (*periphioSpiAccess) createDevice(busNum, chipNum, mode, bits int, maxSpeed int64) (gobot.SpiSystemDevicer, error) {
	return newSpiPeriphIo(busNum, chipNum, mode, bits, maxSpeed)
}
