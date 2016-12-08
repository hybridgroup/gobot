package chip

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)

var _ i2c.I2c = (*Adaptor)(nil)

type NullReadWriteCloser struct {
	contents []byte
}

func (n *NullReadWriteCloser) SetAddress(int) error {
	return nil
}

func (n *NullReadWriteCloser) Write(b []byte) (int, error) {
	n.contents = make([]byte, len(b))
	copy(n.contents[:], b[:])

	return len(b), nil
}

func (n *NullReadWriteCloser) Read(b []byte) (int, error) {
	copy(b, n.contents)
	return len(b), nil
}

var closeErr error

func (n *NullReadWriteCloser) Close() error {
	return closeErr
}

func initTestChipAdaptor() *Adaptor {
	a := NewAdaptor()
	a.Connect()
	return a
}

func TestChipAdaptorDigitalIO(t *testing.T) {
	a := initTestChipAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio1013/value",
		"/sys/class/gpio/gpio1013/direction",
		"/sys/class/gpio/gpio1020/value",
		"/sys/class/gpio/gpio1020/direction",
	})

	sysfs.SetFilesystem(fs)

	a.DigitalWrite("XIO-P0", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio1013/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio1020/value"].Contents = "1"
	i, _ := a.DigitalRead("XIO-P7")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("XIO-P10", 1), errors.New("Not a valid pin"))
}

func TestChipAdaptorI2c(t *testing.T) {
	a := initTestChipAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.I2cStart(0xff)
	a.i2cDevice = &NullReadWriteCloser{}

	a.I2cWrite(0xff, []byte{0x00, 0x01})
	data, _ := a.I2cRead(0xff, 2)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}
