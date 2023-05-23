package system

import (
	"fmt"

	"periph.io/x/conn/v3/physic"
	xspi "periph.io/x/conn/v3/spi"
	xsysfs "periph.io/x/host/v3/sysfs"
)

// spiPeriphIo is the implementation of the SPI interface using the periph.io sysfs implementation for Linux.
type spiPeriphIo struct {
	port xspi.PortCloser
	dev  xspi.Conn
}

// newSpiPeriphIo creates and returns a new connection to a specific SPI device on a bus/chip
// using the periph.io interface.
func newSpiPeriphIo(busNum, chipNum, mode, bits int, maxSpeed int64) (*spiPeriphIo, error) {
	p, err := xsysfs.NewSPI(busNum, chipNum)
	if err != nil {
		return nil, err
	}
	c, err := p.Connect(physic.Frequency(maxSpeed)*physic.Hertz, xspi.Mode(mode), bits)
	if err != nil {
		return nil, err
	}
	return &spiPeriphIo{port: p, dev: c}, nil
}

// TxRx uses the SPI device TX to send/receive data. Implements gobot.SpiSystemDevicer.
func (c *spiPeriphIo) TxRx(tx []byte, rx []byte) error {
	dataLen := len(rx)
	if err := c.dev.Tx(tx, rx); err != nil {
		return err
	}
	if len(rx) != dataLen {
		return fmt.Errorf("Read length (%d) differ to expected (%d)", len(rx), dataLen)
	}
	return nil
}

// Close the SPI connection. Implements gobot.SpiSystemDevicer.
func (c *spiPeriphIo) Close() error {
	return c.port.Close()
}
