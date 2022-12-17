package i2c

import (
	"errors"
	"sync"

	"gobot.io/x/gobot"
)

const (
	// Error event
	Error = "error"
)

const (
	// BusNotInitialized is the initial value for a bus
	BusNotInitialized = -1

	// AddressNotInitialized is the initial value for an address
	AddressNotInitialized = -1
)

var (
	ErrEncryptedBytes  = errors.New("Encrypted bytes")
	ErrNotEnoughBytes  = errors.New("Not enough bytes read")
	ErrNotReady        = errors.New("Device is not ready")
	ErrInvalidPosition = errors.New("Invalid position value")
)

type bitState uint8

const (
	clear bitState = 0x00
	set            = 0x01
)

// Connection is a connection to an I2C device with a specified address
// on a specific bus. Used as an alternative to the I2c interface.
// Implements I2cOperations to talk to the device, wrapping the
// calls in SetAddress to always target the specified device.
// Provided by an Adaptor by implementing the I2cConnector interface.
type Connection gobot.I2cOperations

type i2cConnection struct {
	bus     gobot.I2cSystemDevicer
	address int
	mutex   sync.Mutex
}

// NewConnection creates and returns a new connection to a specific
// i2c device on a bus and address.
func NewConnection(bus gobot.I2cSystemDevicer, address int) (connection *i2cConnection) {
	return &i2cConnection{bus: bus, address: address}
}

// Read data from an i2c device.
func (c *i2cConnection) Read(data []byte) (read int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err = c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	read, err = c.bus.Read(data)
	return
}

// Write data to an i2c device.
func (c *i2cConnection) Write(data []byte) (written int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err = c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	written, err = c.bus.Write(data)
	return
}

// Close connection to i2c device.
func (c *i2cConnection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.bus.Close()
}

// ReadByte reads a single byte from the i2c device.
func (c *i2cConnection) ReadByte() (val byte, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadByte()
}

// ReadByteData reads a byte value for a register on the i2c device.
func (c *i2cConnection) ReadByteData(reg uint8) (val uint8, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadByteData(reg)
}

// ReadWordData reads a word value for a register on the i2c device.
func (c *i2cConnection) ReadWordData(reg uint8) (val uint16, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadWordData(reg)
}

// ReadBlockData reads a block of bytes from a register on the i2c device.
func (c *i2cConnection) ReadBlockData(reg uint8, b []byte) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.ReadBlockData(reg, b)
}

// WriteByte writes a single byte to the i2c device.
func (c *i2cConnection) WriteByte(val byte) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteByte(val)
}

// WriteByteData writes a byte value to a register on the i2c device.
func (c *i2cConnection) WriteByteData(reg uint8, val uint8) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteByteData(reg, val)
}

// WriteWordData writes a word value to a register on the i2c device.
func (c *i2cConnection) WriteWordData(reg uint8, val uint16) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteWordData(reg, val)
}

// WriteBlockData writes a block of bytes to a register on the i2c device.
func (c *i2cConnection) WriteBlockData(reg uint8, b []byte) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteBlockData(reg, b)
}

// WriteBytes writes a block of bytes to the current register on the i2c device.
func (c *i2cConnection) WriteBytes(b []byte) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteBytes(b)
}

// setBit is used to set a bit at a given position to 1.
func setBit(n uint8, pos uint8) uint8 {
	n |= (1 << pos)
	return n
}

// clearBit is used to set a bit at a given position to 0.
func clearBit(n uint8, pos uint8) uint8 {
	mask := ^uint8(1 << pos)
	n &= mask
	return n
}

func twosComplement16Bit(uValue uint16) int16 {
	result := int32(uValue)
	if result&0x8000 != 0 {
		result -= 1 << 16
	}
	return int16(result)
}

func swapBytes(value uint16) uint16 {
	return (value << 8) | (value >> 8)
}
