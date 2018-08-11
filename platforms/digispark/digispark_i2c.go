package digispark

type digisparkI2cConnection struct {
	address uint8
	adaptor *Adaptor
}

// NewDigisparkI2cConnection creates an I2C connection to an I2C device at
// the specified address
func NewDigisparkI2cConnection(adaptor *Adaptor, address uint8) (connection *digisparkI2cConnection) {
	c := &digisparkI2cConnection{adaptor: adaptor, address: address}
	c.Init()
	return c
}

// Init makes sure that the i2c device is already initialized and started
func (c *digisparkI2cConnection) Init() (err error) {
	if !c.adaptor.i2c {
		if err = c.adaptor.littleWire.i2cInit(); err != nil {
			return
		}
		if err = c.adaptor.littleWire.i2cStart(c.address, 0); err != nil { // direction as a param?
			return
		}
		c.adaptor.i2c = true
	}
	return
}

// Read tries to read a full buffer from the i2c device.
// Returns an empty array if the response from the board has timed out.
func (c *digisparkI2cConnection) Read(b []byte) (read int, err error) {
	err = c.Init()
	l := 8
	stop := uint8(0)

	for stop == 0 {
		if read+l >= len(b) {
			l = len(b) - read
			stop = 1
		}
		if err = c.adaptor.littleWire.i2cRead(b[read:read+l], l, stop); err != nil {
			return
		}
		read += l
	}
	return
}

func (c *digisparkI2cConnection) Write(data []byte) (written int, err error) {
	err = c.Init()
	l := 4
	stop := uint8(0)

	for stop == 0 {
		if written+l >= len(data) {
			l = len(data) - written
			stop = 1
		}
		if err = c.adaptor.littleWire.i2cWrite(data[written:written+l], l, stop); err != nil {
			return
		}
		written += l
	}
	return
}

func (c *digisparkI2cConnection) Close() error {
	return nil
}

func (c *digisparkI2cConnection) ReadByte() (val byte, err error) {
	b := make([]byte, 1)
	if err = c.adaptor.littleWire.i2cRead(b, 1, 1); err != nil {
		return
	}
	val = b[0]
	return
}

func (c *digisparkI2cConnection) ReadByteData(reg uint8) (val uint8, err error) {
	if err = c.WriteByte(reg); err != nil {
		return
	}
	return c.ReadByte()
}

func (c *digisparkI2cConnection) ReadWordData(reg uint8) (val uint16, err error) {
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

func (c *digisparkI2cConnection) WriteByte(val byte) (err error) {
	b := []byte{val}
	err = c.adaptor.littleWire.i2cWrite(b, 1, 1)
	return
}

func (c *digisparkI2cConnection) WriteByteData(reg uint8, val byte) (err error) {
	buf := []byte{reg, val}
	_, err = c.Write(buf)
	return
}

func (c *digisparkI2cConnection) WriteWordData(reg uint8, val uint16) (err error) {
	low := uint8(val & 0xff)
	high := uint8((val >> 8) & 0xff)
	buf := []byte{reg, low, high}
	_, err = c.Write(buf)
	return
}

func (c *digisparkI2cConnection) WriteBlockData(reg uint8, data []byte) (err error) {
	if len(data) > 32 {
		data = data[:32]
	}

	buf := make([]byte, len(data)+1)
	copy(buf[:1], data)
	buf[0] = reg
	_, err = c.Write(buf)
	return
}
