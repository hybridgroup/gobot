package spi

import (
	"gobot.io/x/gobot/v2/drivers/common/mfrc522"
)

// MFRC522Driver is a wrapper for SPI bus usage. Please refer to the mfrc522.MFRC522Common package
// for implementation details.
type MFRC522Driver struct {
	*Driver
	*mfrc522.MFRC522Common
}

// NewMFRC522Driver creates a new Gobot Driver for MFRC522 RFID with SPI connection
//
// Params:
//
//	a *Adaptor - the Adaptor to use with this Driver
//
// Optional params:
//
//	 spi.WithBusNumber(int):  bus to use with this driver
//		spi.WithChipNumber(int): chip to use with this driver
//	 spi.WithMode(int):    	 mode to use with this driver
//	 spi.WithBitCount(int):   number of bits to use with this driver
//	 spi.WithSpeed(int64):    speed in Hz to use with this driver
func NewMFRC522Driver(a Connector, options ...func(Config)) *MFRC522Driver {
	d := &MFRC522Driver{
		Driver: NewDriver(a, "MFRC522"),
	}
	d.MFRC522Common = mfrc522.NewMFRC522Common()
	d.afterStart = d.initialize
	for _, option := range options {
		option(d)
	}
	return d
}

func (d *MFRC522Driver) initialize() error {
	wrapper := &conWrapper{origCon: d.connection}
	return d.MFRC522Common.Initialize(wrapper)
}

// this is necessary due to special behavior of shift bytes and set first bit
type conWrapper struct {
	origCon Connection
}

func (w *conWrapper) ReadByteData(reg uint8) (uint8, error) {
	// MSBit=1 for reading, LSBit not used for first byte (address/register)
	return w.origCon.ReadByteData(0x80 | (reg << 1))
}

func (w *conWrapper) WriteByteData(reg uint8, val uint8) error {
	// LSBit not used for first byte (address/register)
	return w.origCon.WriteByteData(reg<<1, val)
}
