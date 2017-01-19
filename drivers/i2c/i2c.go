package i2c

import (
	"errors"

	"gobot.io/x/gobot"
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

type I2cStarter interface {
	// I2cStart initializes I2C communication to a device at a specified address
	// on the default bus for further access using I2cRead() and I2cWrite().
	I2cStart(address int) (err error)
}

type I2cReader interface {
	// I2cStart reads len bytes from a device at the specified address using
	// block reads.
	//
	// Note: There is no way to specify command/register to read from using
	// this interface, it always starts reading from register 0.
	// To read register 42 you have to read 43 bytes (0 to 42)
	// and throw away the first 42 bytes.
	I2cRead(address int, len int) (data []byte, err error)
}

type I2cWriter interface {
	// I2cWrite writes a buffer of bytes to the device at the specified address
	// using block writes.
	// Note that the first byte in the buffer is interpreted as the target
	// command/register.
	I2cWrite(address int, buf []byte) (err error)
}

type I2c interface {
	gobot.Adaptor
	I2cStarter
	I2cReader
	I2cWriter
}

// I2cConnection is a connection to an I2C device with a specified address
// on a specific bus. Used as an alternative to the I2c interface.
// Implements sysfs.SMBusOperations to talk to the device, wrapping the
// calls in SetAddress to always target the specified device.
// Provided by an Adaptor by implementing the I2cConnector interface.
type I2cConnection sysfs.SMBusOperations

type i2cConnection struct {
	bus     sysfs.I2cDevice
	address int
}

func NewI2cConnection(bus sysfs.I2cDevice, address int) (connection *i2cConnection) {
	return &i2cConnection{bus: bus, address: address}
}

func (c *i2cConnection) ReadByte() (val uint8, err error) {
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

func (c *i2cConnection) ReadBlockData(b []byte) (n int, err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return 0, err
	}
	return c.bus.ReadBlockData(b)
}

func (c *i2cConnection) WriteByte(val uint8) (err error) {
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

func (c *i2cConnection) WriteBlockData(b []byte) (err error) {
	if err := c.bus.SetAddress(c.address); err != nil {
		return err
	}
	return c.bus.WriteBlockData(b)
}

// I2cConnector is a replacement to the I2c interface, providing
// access to more than one I2C bus per adaptor, and a more
// fine grained interface to the I2C bus.
type I2cConnector interface {
	I2cGetConnection(address int, bus int) (device I2cConnection, err error)
}
