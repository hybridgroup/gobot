package i2c

import (
	"errors"
	"io"
	"sync"
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

type I2cOperations interface {
	io.ReadWriteCloser
	ReadByte() (val byte, err error)
	ReadByteData(reg uint8) (val uint8, err error)
	ReadWordData(reg uint8) (val uint16, err error)
	WriteByte(val byte) (err error)
	WriteByteData(reg uint8, val uint8) (err error)
	WriteWordData(reg uint8, val uint16) (err error)
	WriteBlockData(reg uint8, b []byte) (err error)
}

// I2cDevice is the interface to a specific i2c bus
type I2cDevice interface {
	I2cOperations
	SetAddress(int) error
}

// Connector lets Adaptors provide the interface for Drivers
// to get access to the I2C buses on platforms that support I2C.
type Connector interface {
	// GetConnection returns a connection to device at the specified address
	// and bus. Bus numbering starts at index 0, the range of valid buses is
	// platform specific.
	GetConnection(address int, bus int) (device Connection, err error)

	// GetDefaultBus returns the default I2C bus index
	GetDefaultBus() int
}

// Connection is a connection to an I2C device with a specified address
// on a specific bus. Used as an alternative to the I2c interface.
// Implements I2cOperations to talk to the device, wrapping the
// calls in SetAddress to always target the specified device.
// Provided by an Adaptor by implementing the I2cConnector interface.
type Connection I2cOperations

type i2cConnection struct {
	bus     I2cDevice
	address int
	mutex   *sync.Mutex
}

// NewConnection creates and returns a new connection to a specific
// i2c device on a bus and address.
func NewConnection(bus I2cDevice, address int) (connection *i2cConnection) {
	return &i2cConnection{bus: bus, address: address, mutex: &sync.Mutex{}}
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
