package i2c

import (
	"errors"
	"fmt"
	"sync"
)

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
	name          string
	bus           int
	address       int
	written       []byte
	mtx           sync.Mutex
	i2cConnectErr bool
	i2cReadImpl   func([]byte) (int, error)
	i2cWriteImpl  func([]byte) (int, error)
}

func (t *i2cTestAdaptor) Testi2cConnectErr(val bool) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.i2cConnectErr = val
}

func (t *i2cTestAdaptor) Testi2cReadImpl(f func([]byte) (int, error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.i2cReadImpl = f
}

func (t *i2cTestAdaptor) Testi2cWriteImpl(f func([]byte) (int, error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.i2cWriteImpl = f
}

func (t *i2cTestAdaptor) Read(b []byte) (count int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.i2cReadImpl(b)
}

func (t *i2cTestAdaptor) Write(b []byte) (count int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.written = append(t.written, b...)
	return t.i2cWriteImpl(b)
}

func (t *i2cTestAdaptor) Close() error {
	return nil
}

func (t *i2cTestAdaptor) ReadByte() (val byte, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	bytes := []byte{0}
	if err = t.readBytes(bytes); err != nil {
		return
	}
	val = bytes[0]
	return
}

func (t *i2cTestAdaptor) ReadByteData(reg uint8) (val uint8, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if err = t.writeBytes([]byte{reg}); err != nil {
		return
	}
	bytes := []byte{0}
	if err = t.readBytes(bytes); err != nil {
		return
	}
	val = bytes[0]
	return
}

func (t *i2cTestAdaptor) ReadWordData(reg uint8) (val uint16, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if err = t.writeBytes([]byte{reg}); err != nil {
		return
	}
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

func (t *i2cTestAdaptor) ReadBlockData(reg uint8, b []byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if err = t.writeBytes([]byte{reg}); err != nil {
		return
	}
	return t.readBytes(b)
}

func (t *i2cTestAdaptor) WriteByte(val byte) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.writeBytes([]byte{val})
}

func (t *i2cTestAdaptor) WriteByteData(reg uint8, val uint8) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	bytes := []byte{reg, val}

	return t.writeBytes(bytes)
}

func (t *i2cTestAdaptor) WriteWordData(reg uint8, val uint16) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	low := uint8(val & 0xff)
	high := uint8((val >> 8) & 0xff)
	bytes := []byte{reg, low, high}

	return t.writeBytes(bytes)
}

func (t *i2cTestAdaptor) WriteBlockData(reg uint8, b []byte) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if len(b) > 32 {
		b = b[:32]
	}
	buf := make([]byte, len(b)+1)
	copy(buf[1:], b)
	buf[0] = reg

	return t.writeBytes(buf)
}

func (t *i2cTestAdaptor) GetConnection(address int, bus int) (connection Connection, err error) {
	if t.i2cConnectErr {
		return nil, errors.New("Invalid i2c connection")
	}
	t.bus = bus
	t.address = address
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
		i2cConnectErr: false,
		i2cReadImpl: func(b []byte) (int, error) {
			return len(b), nil
		},
		i2cWriteImpl: func(b []byte) (int, error) {
			return len(b), nil
		},
	}
}

func (t *i2cTestAdaptor) readBytes(b []byte) error {
	n, err := t.i2cReadImpl(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return fmt.Errorf("Read %v bytes from device by i2c helpers, expected %v", n, len(b))
	}
	return nil
}

func (t *i2cTestAdaptor) writeBytes(b []byte) error {
	t.written = append(t.written, b...)
	// evaluation of count can be done in test
	_, err := t.i2cWriteImpl(b)
	if err != nil {
		return err
	}
	return nil
}
