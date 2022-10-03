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

type bitState uint8

const (
	clear bitState = 0x00
	set            = 0x01
)

// I2cOperations represents the i2c methods according to I2C/SMBus specification.
// Some functions are not in the interface yet:
// * Process Call (WriteWordDataReadWordData)
// * Block Write - Block Read (WriteBlockDataReadBlockData)
// * Host Notify - WriteWordData() can be used instead
//
// see: https://docs.kernel.org/i2c/smbus-protocol.html#key-to-symbols
//
// S: Start condition; Sr: Repeated start condition, used to switch from write to read mode.
// P: Stop condition; Rd/Wr (1 bit): Read/Write bit. Rd equals 1, Wr equals 0.
// A, NA (1 bit): Acknowledge (ACK) and Not Acknowledge (NACK) bit
// Addr (7 bits): I2C 7 bit address. (10 bit I2C address not yet supported by gobot).
// Comm (8 bits): Command byte, a data byte which often selects a register on the device.
// Data (8 bits): A plain data byte. DataLow and DataHigh represent the low and high byte of a 16 bit word.
// Count (8 bits): A data byte containing the length of a block operation.
// [..]: Data sent by I2C device, as opposed to data sent by the host adapter.
//
type I2cOperations interface {
	io.ReadWriteCloser

	// ReadByte must be implemented as the sequence:
	// "S Addr Rd [A] [Data] NA P"
	ReadByte() (val byte, err error)

	// ReadByteData must be implemented as the sequence:
	// "S Addr Wr [A] Comm [A] Sr Addr Rd [A] [Data] NA P"
	ReadByteData(reg uint8) (val uint8, err error)

	// ReadWordData must be implemented as the sequence:
	// "S Addr Wr [A] Comm [A] Sr Addr Rd [A] [DataLow] A [DataHigh] NA P"
	ReadWordData(reg uint8) (val uint16, err error)

	// ReadBlockData must be implemented as the sequence:
	// "S Addr Wr [A] Comm [A] Sr Addr Rd [A] [Count] A [Data] A [Data] A ... A [Data] NA P"
	ReadBlockData(reg uint8, b []byte) (err error)

	// WriteByte must be implemented as the sequence:
	// "S Addr Wr [A] Data [A] P"
	WriteByte(val byte) (err error)

	// WriteByteData must be implemented as the sequence:
	// "S Addr Wr [A] Comm [A] Data [A] P"
	WriteByteData(reg uint8, val uint8) (err error)

	// WriteWordData must be implemented as the sequence:
	// "S Addr Wr [A] Comm [A] DataLow [A] DataHigh [A] P"
	WriteWordData(reg uint8, val uint16) (err error)

	// WriteBlockData must be implemented as the sequence:
	// "S Addr Wr [A] Comm [A] Count [A] Data [A] Data [A] ... [A] Data [A] P"
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
