package spi

import (
	"fmt"
	"sync"

	"gobot.io/x/gobot"
)

// spiConnection is the common implementation of the SPI bus interface.
type spiConnection struct {
	spiSystem gobot.SpiSystemDevicer
	mutex     sync.Mutex
}

// NewConnection uses the given SPI system device and provides it as gobot.SpiOperations
// and Implements gobot.BusOperations.
func NewConnection(spiSystem gobot.SpiSystemDevicer) *spiConnection {
	return &spiConnection{spiSystem: spiSystem}
}

// ReadCommandData uses the SPI device TX to send/receive data. Implements gobot.SpiOperations
func (c *spiConnection) ReadCommandData(command []byte, data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.txRxAndCheckReadLength(command, data)
}

// Close connection to underlying SPI device.
func (c *spiConnection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.spiSystem.Close()
}

// ReadByteData reads a byte from the given register of SPI device. Implements gobot.BusOperations.
func (c *spiConnection) ReadByteData(reg uint8) (uint8, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	buf := []byte{0x0}
	if err := c.txRxAndCheckReadLength([]byte{reg}, buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

// WriteByte writes the given byte value to the current register of SPI device. Implements gobot.BusOperations.
func (c *spiConnection) WriteByte(val byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.writeBytes([]byte{val})
}

// WriteByteData writes the given byte value to the given register of SPI device. Implements gobot.BusOperations.
func (c *spiConnection) WriteByteData(reg byte, data byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.writeBytes([]byte{reg, data})
}

// WriteBlockData writes the given data starting from the given register of SPI device. Implements gobot.BusOperations.
func (c *spiConnection) WriteBlockData(reg byte, data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	buf := make([]byte, len(data)+1)
	copy(buf[1:], data)
	buf[0] = reg
	return c.writeBytes(buf)
}

// WriteBytes writes the given data starting from the current register of bus device. Implements gobot.BusOperations.
func (c *spiConnection) WriteBytes(data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.writeBytes(data)
}

func (c *spiConnection) writeBytes(data []byte) error {
	return c.txRxAndCheckReadLength(data, nil)
}

func (c *spiConnection) txRxAndCheckReadLength(tx []byte, rx []byte) error {
	dataLen := len(rx)
	if err := c.spiSystem.TxRx(tx, rx); err != nil {
		return err
	}
	if len(rx) != dataLen {
		return fmt.Errorf("Read length (%d) differ to expected (%d)", len(rx), dataLen)
	}
	return nil
}
