package digispark

import (
	"errors"
	"fmt"
	"sync"
)

// digisparkI2cConnection implements the interface i2c.I2cOperations
type digisparkI2cConnection struct {
	address uint8
	adaptor *Adaptor
	mtx     sync.Mutex
}

// NewDigisparkI2cConnection creates an i2c connection to an i2c device at
// the specified address
func NewDigisparkI2cConnection(adaptor *Adaptor, address uint8) (connection *digisparkI2cConnection) {
	return &digisparkI2cConnection{adaptor: adaptor, address: address}
}

// Init makes sure that the i2c device is already initialized
func (c *digisparkI2cConnection) Init() (err error) {
	if !c.adaptor.i2c {
		if err = c.adaptor.littleWire.i2cInit(); err != nil {
			return
		}
		c.adaptor.i2c = true
	}
	return
}

// Test tests i2c connection with the given address
func (c *digisparkI2cConnection) Test(address uint8) error {
	if !c.adaptor.i2c {
		return errors.New("Digispark i2c not initialized")
	}
	return c.adaptor.littleWire.i2cStart(address, 0)
}

// UpdateDelay updates i2c signal delay amount; tune if neccessary to fit your requirements
func (c *digisparkI2cConnection) UpdateDelay(duration uint) error {
	if !c.adaptor.i2c {
		return errors.New("Digispark i2c not initialized")
	}
	return c.adaptor.littleWire.i2cUpdateDelay(duration)
}

// Read tries to read a full buffer from the i2c device.
// Returns an empty array if the response from the board has timed out.
func (c *digisparkI2cConnection) Read(b []byte) (countRead int, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.readInternal(b)
}

// Write writes the buffer content in data to the i2c device.
func (c *digisparkI2cConnection) Write(data []byte) (countWritten int, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.writeInternal(data, true)
}

// Close do nothing than return nil.
func (c *digisparkI2cConnection) Close() error {
	return nil
}

// ReadByte reads one byte from the i2c device.
func (c *digisparkI2cConnection) ReadByte() (val byte, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	buf := []byte{0}
	if err = c.readAndCheckCount(buf); err != nil {
		return
	}
	val = buf[0]
	return
}

// ReadByteData reads one byte of the given register address from the i2c device.
func (c *digisparkI2cConnection) ReadByteData(reg uint8) (val uint8, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if err = c.writeAndCheckCount([]byte{reg}, false); err != nil {
		return
	}

	buf := []byte{0}
	if err = c.readAndCheckCount(buf); err != nil {
		return
	}
	val = buf[0]
	return
}

// ReadWordData reads two bytes of the given register address from the i2c device.
func (c *digisparkI2cConnection) ReadWordData(reg uint8) (val uint16, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if err = c.writeAndCheckCount([]byte{reg}, false); err != nil {
		return
	}

	buf := []byte{0, 0}
	if err = c.readAndCheckCount(buf); err != nil {
		return
	}
	low, high := buf[0], buf[1]

	val = (uint16(high) << 8) | uint16(low)
	return
}

// ReadBlockData reads a block of maximum 32 bytes from the given register address of the i2c device.
func (c *digisparkI2cConnection) ReadBlockData(reg uint8, data []byte) (err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if err = c.writeAndCheckCount([]byte{reg}, false); err != nil {
		return
	}

	if len(data) > 32 {
		data = data[:32]
	}
	return c.readAndCheckCount(data)
}

// WriteByte writes one byte to the i2c device.
func (c *digisparkI2cConnection) WriteByte(val byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	buf := []byte{val}
	return c.writeAndCheckCount(buf, true)
}

// WriteByteData writes one byte to the given register address of the i2c device.
func (c *digisparkI2cConnection) WriteByteData(reg uint8, val byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	buf := []byte{reg, val}
	return c.writeAndCheckCount(buf, true)
}

// WriteWordData writes two bytes to the given register address of the i2c device.
func (c *digisparkI2cConnection) WriteWordData(reg uint8, val uint16) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	low := uint8(val & 0xff)
	high := uint8((val >> 8) & 0xff)
	buf := []byte{reg, low, high}
	return c.writeAndCheckCount(buf, true)
}

// WriteBlockData writes a block of maximum 32 bytes to the given register address of the i2c device.
func (c *digisparkI2cConnection) WriteBlockData(reg uint8, data []byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if len(data) > 32 {
		data = data[:32]
	}

	buf := make([]byte, len(data)+1)
	copy(buf[1:], data)
	buf[0] = reg
	return c.writeAndCheckCount(buf, true)
}

func (c *digisparkI2cConnection) readAndCheckCount(buf []byte) error {
	countRead, err := c.readInternal(buf)
	if err != nil {
		return err
	}
	expectedCount := len(buf)
	if countRead != expectedCount {
		return fmt.Errorf("Digispark i2c read %d bytes, expected %d bytes", countRead, expectedCount)
	}
	return nil
}

func (c *digisparkI2cConnection) writeAndCheckCount(buf []byte, finalStop bool) error {
	countWritten, err := c.writeInternal(buf, finalStop)
	if err != nil {
		return err
	}
	expectedCount := len(buf)
	if countWritten != expectedCount {
		return fmt.Errorf("Digispark i2c write %d bytes, expected %d bytes", countWritten, expectedCount)
	}
	return nil
}

func (c *digisparkI2cConnection) readInternal(b []byte) (countRead int, err error) {
	if !c.adaptor.i2c {
		err = errors.New("Digispark i2c not initialized")
		return
	}
	if err = c.adaptor.littleWire.i2cStart(c.address, 1); err != nil {
		return
	}
	l := 8
	stop := uint8(0)

	for stop == 0 {
		if countRead+l >= len(b) {
			l = len(b) - countRead
			stop = 1
		}
		if err = c.adaptor.littleWire.i2cRead(b[countRead:countRead+l], l, stop); err != nil {
			return
		}
		countRead += l
	}
	return
}

func (c *digisparkI2cConnection) writeInternal(data []byte, finalStop bool) (countWritten int, err error) {
	if !c.adaptor.i2c {
		err = errors.New("Digispark i2c not initialized")
		return
	}
	if err = c.adaptor.littleWire.i2cStart(c.address, 0); err != nil {
		return
	}
	l := 4
	lastQuadruplet := false
	stop := uint8(0)

	for !lastQuadruplet {
		if countWritten+l >= len(data) {
			lastQuadruplet = true
			l = len(data) - countWritten
			if finalStop {
				stop = 1
			}
		}
		if err = c.adaptor.littleWire.i2cWrite(data[countWritten:countWritten+l], l, stop); err != nil {
			return
		}
		countWritten += l
	}
	return
}
