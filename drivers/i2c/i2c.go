package i2c

import (
	"errors"

	"gobot.io/x/gobot/sysfs"
)

var (
	ErrEncryptedBytes  = errors.New("Encrypted bytes")
	ErrNotEnoughBytes  = errors.New("Not enough bytes read")
	ErrNotReady        = errors.New("Device is not ready")
	ErrInvalidPosition = errors.New("Invalid position value")
)

const (
	Error = "error"
)

const (
	BusNotInitialized     = -1
	AddressNotInitialized = -1
)

// I2cConnection is a connection to an I2C device with a specified address
// on a specific bus. Used as an alternative to the I2c interface.
// Implements sysfs.I2cOperations to talk to the device, wrapping the
// calls in SetAddress to always target the specified device.
// Provided by an Adaptor by implementing the I2cConnector interface.
type I2cConnection sysfs.I2cOperations

type i2cConnection struct {
	bus     sysfs.I2cDevice
	address int
}

// NewI2cConnection creates and returns a new connection to a specific
// i2c device on a bus and address
func NewI2cConnection(bus sysfs.I2cDevice, address int) (connection *i2cConnection) {
	return &i2cConnection{bus: bus, address: address}
}

// Read data from an i2c device
func (c *i2cConnection) Read(data []byte) (read int, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	read, err = c.bus.Read(data)
	return
}

// Write data to an i2c device
func (c *i2cConnection) Write(data []byte) (written int, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	written, err = c.bus.Write(data)
	return
}

// Close connection to i2c device
func (c *i2cConnection) Close() error {
	return c.bus.Close()
}

// ReadByte reads a single byte from the i2c device
func (c *i2cConnection) ReadByte() (val byte, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadByte()
}

// ReadByteData reads a byte value for a register on the i2c device
func (c *i2cConnection) ReadByteData(reg uint8) (val uint8, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadByteData(reg)
}

// ReadWordData reads a word value for a register on the i2c device
func (c *i2cConnection) ReadWordData(reg uint8) (val uint16, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadWordData(reg)
}

// ReadBlockData reads a block of bytes for a register on the i2c device
func (c *i2cConnection) ReadBlockData(reg uint8, b []byte) (n int, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadBlockData(reg, b)
}

// WriteByte writes a single byte to the i2c device
func (c *i2cConnection) WriteByte(val byte) (err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteByte(val)
}

// WriteByteData writes a byte value to a register on the i2c device
func (c *i2cConnection) WriteByteData(reg uint8, val uint8) (err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteByteData(reg, val)
}

// WriteWordData writes a word value to a register on the i2c device
func (c *i2cConnection) WriteWordData(reg uint8, val uint16) (err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteWordData(reg, val)
}

// WriteBlockData writes a block of bytes to a register on the i2c device
func (c *i2cConnection) WriteBlockData(reg uint8, b []byte) (err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteBlockData(reg, b)
}

// I2cConnector lets Adaotors provide the interface for Drivers
// to get access to the I2C buses on platforms that support I2C.
type I2cConnector interface {
	// I2cGetConnection returns a connection to device at the specified address
	// and bus. Bus numbering starts at index 0, the range of valid buses is
	// platform specific.
	I2cGetConnection(address int, bus int) (device I2cConnection, err error)
	// I2cGetDefaultBus returns the default I2C bus index
	I2cGetDefaultBus() int
}

type i2cConfig struct {
	bus     int
	address int
}

// I2cConfig is the interface which describes how a Driver can specify
// optional I2C params such as which I2C bus it wants to use
type I2cConfig interface {
	// Bus sets which bus to use
	Bus(bus int)

	// GetBus gets which bus to use
	GetBus(def int) int

	// Address sets which address to use
	Address(address int)

	// GetAddress gets which address to use
	GetAddress(def int) int
}

// NewI2cConfig returns a new I2cConfig.
func NewI2cConfig() I2cConfig {
	return &i2cConfig{bus: BusNotInitialized, address: AddressNotInitialized}
}

// Bus sets preferred bus to use
func (i *i2cConfig) Bus(bus int) {
	i.bus = bus
}

// GetBus gets which bus to use
func (i *i2cConfig) GetBus(d int) int {
	if i.bus == BusNotInitialized {
		return d
	}

	return i.bus
}

// Bus sets which bus to use as a optional param
func Bus(bus int) func(I2cConfig) {
	return func(i I2cConfig) {
		i.Bus(bus)
	}
}

// Address sets which address to use
func (i *i2cConfig) Address(address int) {
	i.address = address
}

// GetAddress gets which address to use
func (i *i2cConfig) GetAddress(a int) int {
	if i.address == AddressNotInitialized {
		return a
	}

	return i.address
}

// Address sets which address to use as a optional param
func Address(address int) func(I2cConfig) {
	return func(i I2cConfig) {
		i.Address(address)
	}
}
