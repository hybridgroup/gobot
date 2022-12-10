package system

import (
	"periph.io/x/conn/v3/physic"
	xspi "periph.io/x/conn/v3/spi"
	xsysfs "periph.io/x/host/v3/sysfs"
)

// spiConnectionPeriphIo is the implementation of the SPI interface using the periph.io
// sysfs implementation for Linux.
type spiConnectionPeriphIo struct {
	port xspi.PortCloser
	dev  xspi.Conn
}

// NewspiConnectionPeriphIo creates and returns a new connection to a specific
// spi device on a bus/chip using the periph.io interface.
func (a *Accesser) NewSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (*spiConnectionPeriphIo, error) {
	p, err := xsysfs.NewSPI(busNum, chipNum)
	if err != nil {
		return nil, err
	}
	c, err := p.Connect(physic.Frequency(maxSpeed)*physic.Hertz, xspi.Mode(mode), bits)
	if err != nil {
		return nil, err
	}
	return &spiConnectionPeriphIo{port: p, dev: c}, nil
}

// Close the SPI connection.
func (c *spiConnectionPeriphIo) Close() error {
	return c.port.Close()
}

// Tx uses the SPI device to send/receive data.
func (c *spiConnectionPeriphIo) Tx(w, r []byte) error {
	return c.dev.Tx(w, r)
}
