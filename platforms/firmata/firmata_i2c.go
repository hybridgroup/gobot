package firmata

import (
	//	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata/client"
)

type firmataI2cConnection struct {
	address int
	adaptor *Adaptor
}

// NewFirmataI2cConnection creates an I2C connection to an I2C device at
// the specified address
func NewFirmataI2cConnection(adaptor *Adaptor, address int) (connection *firmataI2cConnection) {
	return &firmataI2cConnection{adaptor: adaptor, address: address}
}

// Read tries to read a full buffer from the i2c device.
// Returns an empty array if the response from the board has timed out.
func (c *firmataI2cConnection) Read(b []byte) (read int, err error) {
	ret := make(chan []byte)

	if err = c.adaptor.Board.I2cRead(c.address, len(b)); err != nil {
		return
	}

	c.adaptor.Board.Once(c.adaptor.Board.Event("I2cReply"), func(data interface{}) {
		ret <- data.(client.I2cReply).Data
	})

	result := <-ret
	copy(b, result)

	read = len(result)

	return
}

func (c *firmataI2cConnection) Write(data []byte) (written int, err error) {
	var chunk []byte
	for len(data) >= 16 {
		chunk, data = data[:16], data[16:]
		err = c.adaptor.Board.I2cWrite(c.address, chunk)
		if err != nil {
			return
		}
		written += len(chunk)
	}
	if len(data) > 0 {
		err = c.adaptor.Board.I2cWrite(c.address, data[:])
		written += len(data)
	}
	return
}

func (c *firmataI2cConnection) Close() error {
	return nil
}

func (c *firmataI2cConnection) ReadByte() (val byte, err error) {
	buf := []byte{0}
	if _, err = c.Read(buf); err != nil {
		return
	}
	val = buf[0]
	return
}

func (c *firmataI2cConnection) ReadByteData(reg uint8) (val uint8, err error) {
	if err = c.WriteByte(reg); err != nil {
		return
	}
	return c.ReadByte()
}

func (c *firmataI2cConnection) ReadWordData(reg uint8) (val uint16, err error) {
	if err = c.WriteByte(reg); err != nil {
		return
	}

	buf := []byte{0, 0}
	if _, err = c.Read(buf); err != nil {
		return
	}
	low, high := buf[0], buf[1]

	val = (uint16(high) << 8) | uint16(low)
	return
}

func (c *firmataI2cConnection) WriteByte(val byte) (err error) {
	buf := []byte{val}
	_, err = c.Write(buf)
	return
}

func (c *firmataI2cConnection) WriteByteData(reg uint8, val byte) (err error) {
	buf := []byte{reg, val}
	_, err = c.Write(buf)
	return
}

func (c *firmataI2cConnection) WriteWordData(reg uint8, val uint16) (err error) {
	low := uint8(val & 0xff)
	high := uint8((val >> 8) & 0xff)
	buf := []byte{reg, low, high}
	_, err = c.Write(buf)
	return
}

func (c *firmataI2cConnection) WriteBlockData(reg uint8, data []byte) (err error) {
	if len(data) > 32 {
		data = data[:32]
	}

	buf := make([]byte, len(data)+1)
	copy(buf[:1], data)
	buf[0] = reg
	_, err = c.Write(buf)
	return
}
