package spi

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/system"
)

// make sure that this SpiBusAdaptor fulfills all the required interfaces
var _ Connector = (*spiTestAdaptor)(nil)
var _ gobot.Connection = (*spiTestAdaptor)(nil)

type spiTestAdaptor struct {
	sys *system.Accesser
	//busNum        int
	spiConnectErr bool
	spi           *system.MockSpiAccess
	connection    Connection
}

func newSpiTestAdaptor() *spiTestAdaptor {
	sys := system.NewAccesser()
	spi := sys.UseMockSpi()
	a := &spiTestAdaptor{
		sys: sys,
		spi: spi,
	}
	return a
}

// spi.Connector interfaces
func (a *spiTestAdaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (Connection, error) {
	if a.spiConnectErr {
		return nil, fmt.Errorf("Invalid SPI connection in helper")
	}
	//a.busNum = busNum
	con, err := a.sys.NewSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	a.connection = con
	return con, err
}

func (a *spiTestAdaptor) SpiDefaultBusNumber() int  { return 0 }
func (a *spiTestAdaptor) SpiDefaultChipNumber() int { return 0 }
func (a *spiTestAdaptor) SpiDefaultMode() int       { return 0 }
func (a *spiTestAdaptor) SpiDefaultBitCount() int   { return 0 }
func (a *spiTestAdaptor) SpiDefaultMaxSpeed() int64 { return 0 }

// gobot.Connection interfaces
func (a *spiTestAdaptor) Connect() error  { return nil }
func (a *spiTestAdaptor) Finalize() error { return nil }
func (a *spiTestAdaptor) Name() string    { return "board name" }
func (a *spiTestAdaptor) SetName(string)  { return }
