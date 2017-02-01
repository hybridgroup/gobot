package i2c

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
	pos          int
	i2cReadImpl  func() ([]byte, error)
	i2cWriteImpl func() error
	i2cStartImpl func() error
}

func (t *i2cTestAdaptor) I2cStart(int) (err error) {
	t.pos = 0
	return t.i2cStartImpl()
}

func (t *i2cTestAdaptor) I2cRead(int, int) (data []byte, err error) {
	return t.i2cReadImpl()
}

func (t *i2cTestAdaptor) I2cWrite(address int, b []byte) (err error) {
	t.written = append(t.written, b...)
	return t.i2cWriteImpl()
}

func (t *i2cTestAdaptor) ReadByte() (val uint8, err error) {
	bytes, err := t.i2cReadImpl()
	val = bytes[t.pos]
	t.pos++
	return
}

func (t *i2cTestAdaptor) ReadByteData(reg uint8) (val uint8, err error) {
	bytes, err := t.i2cReadImpl()
	return bytes[reg], err
}

func (t *i2cTestAdaptor) ReadWordData(reg uint8) (val uint16, err error) {
	bytes, err := t.i2cReadImpl()
	low, high := bytes[reg], bytes[reg + 1]
	return (uint16(high) << 8) | uint16(low), err
}

func (t *i2cTestAdaptor) ReadBlockData(b []byte) (n int, err error) {
	reg := b[0]
	bytes, err := t.i2cReadImpl()
	copy(b, bytes[reg:])
	return len(b), err
}

func (t *i2cTestAdaptor) WriteByte(val uint8) (err error) {
	t.pos = int(val)
	t.written = append(t.written, val)
	return t.i2cWriteImpl()
}

func (t *i2cTestAdaptor) WriteByteData(reg uint8, val uint8) (err error) {
	t.pos = int(reg)
	t.written = append(t.written, reg)
	t.written = append(t.written, val)
	return t.i2cWriteImpl()
}

func (t *i2cTestAdaptor) WriteBlockData(b []byte) (err error) {
	t.pos = int(b[0])
	t.written = append(t.written, b...)
	return t.i2cWriteImpl()
}

func (t *i2cTestAdaptor) I2cGetConnection( /* address */ int /* bus */, int) (connection I2cConnection, err error) {
	return t, nil
}

func (t *i2cTestAdaptor) Name() string          { return t.name }
func (t *i2cTestAdaptor) SetName(n string)      { t.name = n }
func (t *i2cTestAdaptor) Connect() (err error)  { return }
func (t *i2cTestAdaptor) Finalize() (err error) { return }

func newI2cTestAdaptor() *i2cTestAdaptor {
	return &i2cTestAdaptor{
		i2cReadImpl: func() ([]byte, error) {
			return []byte{}, nil
		},
		i2cWriteImpl: func() error {
			return nil
		},
		i2cStartImpl: func() error {
			return nil
		},
	}
}
