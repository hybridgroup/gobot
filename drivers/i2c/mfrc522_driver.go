package i2c

import (
	"gobot.io/x/gobot/v2/drivers/common/mfrc522"
)

const mfrc522DefaultAddress = 0x00

// MFRC522Driver is a wrapper for i2c bus usage. Please refer to the mfrc522.MFRC522Common package
// for implementation details.
type MFRC522Driver struct {
	*Driver
	*mfrc522.MFRC522Common
}

// NewMFRC522Driver creates a new Gobot Driver for MFRC522 RFID with i2c connection
//
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewMFRC522Driver(c Connector, options ...func(Config)) *MFRC522Driver {
	d := &MFRC522Driver{
		Driver: NewDriver(c, "MFRC522", mfrc522DefaultAddress),
	}
	d.MFRC522Common = mfrc522.NewMFRC522Common()
	d.afterStart = d.initialize
	for _, option := range options {
		option(d)
	}
	return d
}

func (d *MFRC522Driver) initialize() error {
	return d.MFRC522Common.Initialize(d.connection)
}
