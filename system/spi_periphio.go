package system

import (
	"fmt"

	"periph.io/x/conn/v3/physic"
	xspi "periph.io/x/conn/v3/spi"
	xsysfs "periph.io/x/host/v3/sysfs"
)

// spiConnectionPeriphIo is the implementation of the SPI interface using the periph.io sysfs implementation for Linux.
type spiConnectionPeriphIo struct {
	port xspi.PortCloser
	dev  xspi.Conn
}

// newSpiConnectionPeriphIo creates and returns a new connection to a specific SPI device on a bus/chip
// using the periph.io interface.
func newSpiConnectionPeriphIo(busNum, chipNum, mode, bits int, maxSpeed int64) (*spiConnectionPeriphIo, error) {
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

// ReadData uses the SPI device TX to send/receive data.
func (c *spiConnectionPeriphIo) ReadData(command []byte, data []byte) error {
	dataLen := len(data)
	if err := c.dev.Tx(command, data); err != nil {
		return err
	}
	if len(data) != dataLen {
		return fmt.Errorf("Read length (%d) differ to expected (%d)", len(data), dataLen)
	}
	return nil
}

// WriteData uses the SPI device TX to send data.
func (c *spiConnectionPeriphIo) WriteData(data []byte) error {
	return c.dev.Tx(data, nil)
}
