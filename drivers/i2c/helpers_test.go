package i2c

import "fmt"

var rgb = map[string]interface{}{
	"red":   1.0,
	"green": 1.0,
	"blue":  1.0,
}

func castColor(color string) byte {
	return byte(rgb[color].(float64))
}

var red = castColor("red")
var green = castColor("green")
var blue = castColor("blue")

type i2cTestAdaptor struct {
	name         string
	written      []byte
	i2cReadImpl  func([]byte) (int, error)
	i2cWriteImpl func([]byte) (int, error)
}

func (t *i2cTestAdaptor) Read(b []byte) (count int, err error) {
	return t.i2cReadImpl(b)
}

func (t *i2cTestAdaptor) Write(b []byte) (count int, err error) {
	t.written = append(t.written, b...)
	return t.i2cWriteImpl(b)
}

func (t *i2cTestAdaptor) Close() error {
	return nil
}

func (t *i2cTestAdaptor) ReadByte() (val byte, err error) {
	bytes := []byte{0}
	bytesRead, err := t.i2cReadImpl(bytes)
	if err != nil {
		return 0, err
	}
	if bytesRead != 1 {
		return 0, fmt.Errorf("Buffer underrun")
	}
	val = bytes[0]
	return
}

func (t *i2cTestAdaptor) ReadByteData(reg uint8) (val uint8, err error) {
	bytes := []byte{0}
	bytesRead, err := t.i2cReadImpl(bytes)
	if err != nil {
		return 0, err
	}
	if bytesRead != 1 {
		return 0, fmt.Errorf("Buffer underrun")
	}
	val = bytes[0]
	return
}

func (t *i2cTestAdaptor) ReadWordData(reg uint8) (val uint16, err error) {
	bytes := []byte{0, 0}
	bytesRead, err := t.i2cReadImpl(bytes)
	if err != nil {
		return 0, err
	}
	if bytesRead != 2 {
		return 0, fmt.Errorf("Buffer underrun")
	}
	low, high := bytes[0], bytes[1]
	return (uint16(high) << 8) | uint16(low), err
}

func (t *i2cTestAdaptor) ReadBlockData(_ uint8, b []byte) (n int, err error) {
	bytes := make([]byte, 32)
	bytesRead, err := t.i2cReadImpl(bytes)
	copy(b, bytes[:bytesRead])
	return bytesRead, err
}

func (t *i2cTestAdaptor) WriteByte(val byte) (err error) {
	t.written = append(t.written, val)
	bytes := []byte{val}
	_, err = t.i2cWriteImpl(bytes)
	return
}

func (t *i2cTestAdaptor) WriteByteData(reg uint8, val uint8) (err error) {
	t.written = append(t.written, reg)
	t.written = append(t.written, val)
	bytes := []byte{val}
	_, err = t.i2cWriteImpl(bytes)
	return
}

func (t *i2cTestAdaptor) WriteWordData(reg uint8, val uint16) (err error) {
	t.written = append(t.written, reg)
	low := uint8(val & 0xff)
	high := uint8((val >> 8) & 0xff)
	t.written = append(t.written, low)
	t.written = append(t.written, high)
	bytes := []byte{low, high}
	_, err = t.i2cWriteImpl(bytes)
	return
}

func (t *i2cTestAdaptor) WriteBlockData(reg uint8, b []byte) (err error) {
	t.written = append(t.written, reg)
	t.written = append(t.written, b...)
	_, err = t.i2cWriteImpl(b)
	return
}

func (t *i2cTestAdaptor) GetConnection( /* address */ int /* bus */, int) (connection Connection, err error) {
	return t, nil
}

func (t *i2cTestAdaptor) GetDefaultBus() int {
	return 0
}

func (t *i2cTestAdaptor) Name() string          { return t.name }
func (t *i2cTestAdaptor) SetName(n string)      { t.name = n }
func (t *i2cTestAdaptor) Connect() (err error)  { return }
func (t *i2cTestAdaptor) Finalize() (err error) { return }

func newI2cTestAdaptor() *i2cTestAdaptor {
	return &i2cTestAdaptor{
		i2cReadImpl: func([]byte) (int, error) {
			return 0, nil
		},
		i2cWriteImpl: func([]byte) (int, error) {
			return 0, nil
		},
	}
}
