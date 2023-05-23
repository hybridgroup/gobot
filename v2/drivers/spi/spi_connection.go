package spi

import (
	"fmt"
	"sync"

	"gobot.io/x/gobot/v2"
)

const (
	spiDebugByte  = false
	spiDebugBlock = false
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
// On write command, the first byte normally contains the address and mode.
// On read data, the return value is most likely one byte behind the command.
// The length of command and data needs to be the same (except data is nil).
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
	if err := c.readAlignedBlockData(reg, buf); err != nil {
		return 0, err
	}

	if spiDebugByte {
		fmt.Printf("ReadByteData: register 0x%02X/0x%02X : 0x%02X %dd\n", reg, reg&0x7F>>1, buf[0], buf[0])
	}
	return buf[0], nil
}

// ReadBlockData fills the given buffer with reads starting from the given register of SPI device.
// Implements gobot.BusOperations.
func (c *spiConnection) ReadBlockData(reg uint8, data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.readAlignedBlockData(reg, data); err != nil {
		return err
	}

	if spiDebugBlock {
		fmt.Printf("ReadBlockData: register 0x%02X/0x%02X : %v\n", reg, reg&0x7F>>1, data)
	}
	return nil
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

func (c *spiConnection) readAlignedBlockData(reg uint8, data []byte) error {
	// length of TX needs to equal length of RX
	// the read value is one cycle behind the write, so for n bytes to read, we need n+1 bytes (to read and write)
	buflen := len(data) + 1
	writeBuf := make([]byte, buflen)
	readBuf := make([]byte, buflen)
	writeBuf[0] = reg
	if err := c.txRxAndCheckReadLength(writeBuf, readBuf); err != nil {
		return err
	}
	copy(data, readBuf[1:])
	return nil
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
