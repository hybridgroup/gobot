package system

import (
	"periph.io/x/conn/v3/physic"
	xspi "periph.io/x/conn/v3/spi"
	xsysfs "periph.io/x/host/v3/sysfs"
)

// SpiConnection is the implementation of the SPI interface using the periph.io
// sysfs implementation for Linux.
type SpiConnection struct {
	port xspi.PortCloser
	dev  xspi.Conn
}

// NewSpiConnection creates and returns a new connection to a specific
// spi device on a bus/chip using the periph.io interface.
func (a *Accesser) NewSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (*SpiConnection, error) {
	p, err := xsysfs.NewSPI(busNum, chipNum)
	if err != nil {
		return nil, err
	}
	c, err := p.Connect(physic.Frequency(maxSpeed)*physic.Hertz, xspi.Mode(mode), bits)
	if err != nil {
		return nil, err
	}
	return &SpiConnection{port: p, dev: c}, nil
}

// Close the SPI connection.
func (c *SpiConnection) Close() error {
	return c.port.Close()
}

// Tx uses the SPI device to send/receive data.
func (c *SpiConnection) Tx(w, r []byte) error {
	return c.dev.Tx(w, r)
}
