package spi

import (
	"time"

	xspi "golang.org/x/exp/io/spi"
)

type TestConnector struct{}

func (ctr *TestConnector) GetSpiConnection(busNum, mode int, maxSpeed int64) (device Connection, err error) {
	return NewConnection(&TestSpiDevice{}), nil
}

func (ctr *TestConnector) GetSpiDefaultBus() int {
	return 0
}

func (ctr *TestConnector) GetSpiDefaultMode() int {
	return 0
}

func (ctr *TestConnector) GetSpiDefaultMaxSpeed() int64 {
	return 0
}

type TestSpiDevice struct {
	bus SPIDevice
}

func (c *TestSpiDevice) Close() error {
	return nil
}

func (c *TestSpiDevice) SetBitOrder(o xspi.Order) error {
	return nil
}

func (c *TestSpiDevice) SetBitsPerWord(bits int) error {
	return nil
}

func (c *TestSpiDevice) SetCSChange(leaveEnabled bool) error {
	return nil
}

func (c *TestSpiDevice) SetDelay(t time.Duration) error {
	return nil
}

func (c *TestSpiDevice) SetMaxSpeed(speed int) error {
	return nil
}

func (c *TestSpiDevice) SetMode(mode xspi.Mode) error {
	return nil
}

func (c *TestSpiDevice) Tx(w, r []byte) error {
	return nil
}
