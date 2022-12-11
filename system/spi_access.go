package system

import "gobot.io/x/gobot"

type periphioSpiAccess struct{}

func (*periphioSpiAccess) createConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (gobot.SpiOperations, error) {
	return newSpiConnectionPeriphIo(busNum, chipNum, mode, bits, maxSpeed)
}
