package i2c

import (
	"errors"

	"gobot.io/x/gobot/sysfs"
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

// Connection is a connection to an I2C device with a specified address
// on a specific bus. Used as an alternative to the I2c interface.
// Implements sysfs.I2cOperations to talk to the device, wrapping the
// calls in SetAddress to always target the specified device.
// Provided by an Adaptor by implementing the I2cConnector interface.
type Connection sysfs.I2cOperations

type i2cConnection struct {
	bus     sysfs.I2cDevice
	address int
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

type i2cConfig struct {
	bus     int
	address int
}

// Config is the interface which describes how a Driver can specify
// optional I2C params such as which I2C bus it wants to use.
type Config interface {
	// WithBus sets which bus to use
	WithBus(bus int)

	// GetBusOrDefault gets which bus to use
	GetBusOrDefault(def int) int

	// WithAddress sets which address to use
	WithAddress(address int)

	// GetAddressOrDefault gets which address to use
	GetAddressOrDefault(def int) int
}

// NewConnection creates and returns a new connection to a specific
// i2c device on a bus and address.
func NewConnection(bus sysfs.I2cDevice, address int) (connection *i2cConnection) {
	return &i2cConnection{bus: bus, address: address}
}

// NewConfig returns a new I2c Config.
func NewConfig() Config {
	return &i2cConfig{bus: BusNotInitialized, address: AddressNotInitialized}
}

// Read data from an i2c device.
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

// WithBus sets preferred bus to use.
func (i *i2cConfig) WithBus(bus int) {
	i.bus = bus
}

// GetBusOrDefault returns which bus to use, either the one set using WithBus(),
// or the default value which is passed in as the one param.
func (i *i2cConfig) GetBusOrDefault(d int) int {
	if i.bus == BusNotInitialized {
		return d
	}

	return i.bus
}

// WithBus sets which bus to use as a optional param
func WithBus(bus int) func(Config) {
	return func(i Config) {
		i.WithBus(bus)
	}
}

// WithAddress sets which address to use
func (i *i2cConfig) WithAddress(address int) {
	i.address = address
}

// GetAddressOrDefault returns which address to use, either
// the one set using WithBus(), or the default value which
// is passed in as the param.
func (i *i2cConfig) GetAddressOrDefault(a int) int {
	if i.address == AddressNotInitialized {
		return a
	}

	return i.address
}

// WithAddress sets which address to use as a optional param
func WithAddress(address int) func(Config) {
	return func(i Config) {
		i.WithAddress(address)
	}
}
