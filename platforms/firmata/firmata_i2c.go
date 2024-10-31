//go:build !windows
// +build !windows

package firmata

import (
	"fmt"
	"sync"

	"gobot.io/x/gobot/v2/platforms/firmata/client"
)

// firmataI2cConnection implements the interface gobot.I2cOperations
type firmataI2cConnection struct {
	address int
	adaptor *Adaptor
	mtx     sync.Mutex
}

// NewFirmataI2cConnection creates an I2C connection to an I2C device at
// the specified address
func NewFirmataI2cConnection(adaptor *Adaptor, address int) *firmataI2cConnection {
	return &firmataI2cConnection{adaptor: adaptor, address: address}
}

// Read tries to read a full buffer from the i2c device.
// Returns an empty array if the response from the board has timed out.
func (c *firmataI2cConnection) Read(b []byte) (int, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.readInternal(b)
}

// Write writes the buffer content in data to the i2c device.
func (c *firmataI2cConnection) Write(data []byte) (int, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.writeInternal(data)
}

// Close do nothing than return nil.
func (c *firmataI2cConnection) Close() error {
	return nil
}

// ReadByte reads one byte from the i2c device.
func (c *firmataI2cConnection) ReadByte() (byte, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	buf := []byte{0}
	if err := c.readAndCheckCount(buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

// ReadByteData reads one byte of the given register address from the i2c device.
// TODO: implement the specification, because some devices will not work with this
//
//	current:  "S Addr Wr [A] Comm [A] P S Addr Rd [A] [Data] NA P"
//	required: "S Addr Wr [A] Comm [A] S Addr Rd [A] [Data] NA P"
func (c *firmataI2cConnection) ReadByteData(reg uint8) (uint8, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if err := c.writeAndCheckCount([]byte{reg}); err != nil {
		return 0, err
	}

	buf := []byte{0}
	if err := c.readAndCheckCount(buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

// ReadWordData reads two bytes of the given register address from the i2c device.
// TODO: implement the specification, because some devices will not work with this
//
//	current:  "S Addr Wr [A] Comm [A] P S Addr Rd [A] [DataLow] A [DataHigh] NA P"
//	required: "S Addr Wr [A] Comm [A] S Addr Rd [A] [DataLow] A [DataHigh] NA P"
func (c *firmataI2cConnection) ReadWordData(reg uint8) (uint16, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if err := c.writeAndCheckCount([]byte{reg}); err != nil {
		return uint16(0), err
	}

	buf := []byte{0, 0}
	if err := c.readAndCheckCount(buf); err != nil {
		return uint16(0), err
	}
	low, high := buf[0], buf[1]
	return (uint16(high) << 8) | uint16(low), nil
}

// ReadBlockData reads a block of maximum 32 bytes from the given register address of the i2c device.
// TODO: implement the specification, because some devices will not work with this
//
//	current:  "S Addr Wr [A] Comm [A] P S Addr Rd [A] [Count] A [Data] A [Data] A ... A [Data] NA P"
//	required: "S Addr Wr [A] Comm [A] S Addr Rd [A] [Count] A [Data] A [Data] A ... A [Data] NA P"
func (c *firmataI2cConnection) ReadBlockData(reg uint8, data []byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if err := c.writeAndCheckCount([]byte{reg}); err != nil {
		return err
	}

	if len(data) > 32 {
		data = data[:32]
	}
	return c.readAndCheckCount(data)
}

// WriteByte writes one byte to the i2c device.
func (c *firmataI2cConnection) WriteByte(val byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	buf := []byte{val}
	return c.writeAndCheckCount(buf)
}

// WriteByteData writes one byte to the given register address of the i2c device.
func (c *firmataI2cConnection) WriteByteData(reg uint8, val byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	buf := []byte{reg, val}
	return c.writeAndCheckCount(buf)
}

// WriteWordData writes two bytes to the given register address of the i2c device.
func (c *firmataI2cConnection) WriteWordData(reg uint8, val uint16) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	low := uint8(val & 0xff)         //nolint:gosec // ok here
	high := uint8((val >> 8) & 0xff) //nolint:gosec // ok here
	buf := []byte{reg, low, high}
	return c.writeAndCheckCount(buf)
}

// WriteBlockData writes a block of maximum 32 bytes to the given register address of the i2c device.
func (c *firmataI2cConnection) WriteBlockData(reg uint8, data []byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if len(data) > 32 {
		data = data[:32]
	}

	buf := make([]byte, len(data)+1)
	copy(buf[1:], data)
	buf[0] = reg
	return c.writeAndCheckCount(buf)
}

// WriteBytes writes a block of maximum 32 bytes to the current register address of the i2c device.
func (c *firmataI2cConnection) WriteBytes(buf []byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if len(buf) > 32 {
		buf = buf[:32]
	}

	return c.writeAndCheckCount(buf)
}

func (c *firmataI2cConnection) readAndCheckCount(buf []byte) error {
	countRead, err := c.readInternal(buf)
	if err != nil {
		return err
	}
	expectedCount := len(buf)
	if countRead != expectedCount {
		return fmt.Errorf("Firmata i2c read %d bytes, expected %d bytes", countRead, expectedCount)
	}
	return nil
}

func (c *firmataI2cConnection) writeAndCheckCount(buf []byte) error {
	countWritten, err := c.writeInternal(buf)
	if err != nil {
		return err
	}
	expectedCount := len(buf)
	if countWritten != expectedCount {
		return fmt.Errorf("Firmata i2c write %d bytes, expected %d bytes", countWritten, expectedCount)
	}
	return nil
}

func (c *firmataI2cConnection) readInternal(b []byte) (int, error) {
	ret := make(chan []byte)

	if err := c.adaptor.Board.I2cRead(c.address, len(b)); err != nil {
		return 0, err
	}

	if err := c.adaptor.Board.Once(c.adaptor.Board.Event("I2cReply"), func(data interface{}) {
		ret <- data.(client.I2cReply).Data //nolint:forcetypeassert // ok here
	}); err != nil {
		return 0, err
	}

	result := <-ret
	copy(b, result)

	return len(result), nil
}

func (c *firmataI2cConnection) writeInternal(data []byte) (int, error) {
	var chunk []byte
	var written int
	for len(data) >= 16 {
		chunk, data = data[:16], data[16:]
		if err := c.adaptor.Board.I2cWrite(c.address, chunk); err != nil {
			return written, err
		}
		written += len(chunk)
	}
	if len(data) > 0 {
		if err := c.adaptor.Board.I2cWrite(c.address, data); err != nil {
			return written, err
		}
		written += len(data)
	}
	return written, nil
}
