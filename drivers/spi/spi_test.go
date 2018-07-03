package spi

import (
	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/physic"
	xspi "periph.io/x/periph/conn/spi"
)

type TestConnector struct{}

func (ctr *TestConnector) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (device Connection, err error) {
	return NewConnection(&TestSpiConnection{}, &TestSpiDevice{}), nil
}

func (ctr *TestConnector) GetSpiDefaultBus() int {
	return 0
}

func (ctr *TestConnector) GetSpiDefaultChip() int {
	return 0
}

func (ctr *TestConnector) GetSpiDefaultMode() int {
	return 0
}

func (ctr *TestConnector) GetSpiDefaultBits() int {
	return 0
}

func (ctr *TestConnector) GetSpiDefaultMaxSpeed() int64 {
	return 0
}

type TestSpiDevice struct {
	dev Connection
}

func (c *TestSpiDevice) String() string {
	return "TestSpiDevice"
}

func (c *TestSpiDevice) Duplex() conn.Duplex {
	return conn.Half
}

func (c *TestSpiDevice) TxPackets(p []xspi.Packet) error {
	return nil
}

func (c *TestSpiDevice) Tx(w, r []byte) error {
	return nil
}

type TestSpiConnection struct {
	conn Operations
}

func (c *TestSpiConnection) String() string {
	return "TestSpiConnection"
}

func (c *TestSpiConnection) Close() error {
	return nil
}

func (c *TestSpiConnection) Connect(maxHz physic.Frequency, mode xspi.Mode, bits int) (xspi.Conn, error) {
	return nil, nil
}

func (c *TestSpiConnection) LimitSpeed(maxHz physic.Frequency) error {
	return nil
}
