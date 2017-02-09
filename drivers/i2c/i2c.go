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
	Error    = "error"
	Joystick = "joystick"
	C        = "c"
	Z        = "z"
)

const (
	BusNotInitialized = -1
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

func NewI2cConnection(bus sysfs.I2cDevice, address int) (connection *i2cConnection) {
	return &i2cConnection{bus: bus, address: address}
}

func (c *i2cConnection) Read(data []byte) (read int, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	read, err = c.bus.Read(data)
	return
}

func (c *i2cConnection) Write(data []byte) (written int, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	written, err = c.bus.Write(data)
	return
}

func (c *i2cConnection) Close() error {
	return c.bus.Close()
}

func (c *i2cConnection) ReadByte() (val byte, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadByte()
}

func (c *i2cConnection) ReadByteData(reg uint8) (val uint8, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadByteData(reg)
}

func (c *i2cConnection) ReadWordData(reg uint8) (val uint16, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadWordData(reg)
}

func (c *i2cConnection) ReadBlockData(reg uint8, b []byte) (n int, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadBlockData(reg, b)
}

func (c *i2cConnection) WriteByte(val byte) (err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteByte(val)
}

func (c *i2cConnection) WriteByteData(reg uint8, val uint8) (err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteByteData(reg, val)
}

func (c *i2cConnection) WriteWordData(reg uint8, val uint16) (err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteWordData(reg, val)
}

func (c *i2cConnection) WriteBlockData(reg uint8, b []byte) (err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteBlockData(reg, b)
}

// I2cConnector provides access to the I2C buses on platforms that support them.
type I2cConnector interface {
	// I2cGetConnection returns a connection to device at the specified address
	// and bus. Bus numbering starts at index 0, the range of valid buses is
	// platform specific.
	I2cGetConnection(address int, bus int) (device I2cConnection, err error)
	// I2cGetDefaultBus returns the default I2C bus index
	I2cGetDefaultBus() int
}

type i2cBusser struct {
	bus int
}

// I2cBusser is the interface which describes how a Driver can specify
// which I2C bus it wants to use
type I2cBusser interface {
	// Bus sets which bus to use
	Bus(bus int)

	// GetBus gets which bus to use
	GetBus() int
}

// NewI2cBusser returns a new I2cBusser.
func NewI2cBusser() I2cBusser {
	return &i2cBusser{}
}

// Bus sets which bus to use
func (i *i2cBusser) Bus(bus int) {
	i.bus = bus
}

// GetBus gets which bus to use
func (i *i2cBusser) GetBus() int {
	return i.bus
}

// Bus sets which bus to use as a optional param
func Bus(bus int) func(I2cBusser) {
	return func(i I2cBusser) {
		i.Bus(bus)
	}
}
