package chip

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

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

var closeErr error = nil

func (n *NullReadWriteCloser) Close() error {
	return closeErr
}

func initTestChipAdaptor() *ChipAdaptor {
	a := NewChipAdaptor("myAdaptor")
	a.Connect()
	return a
}

func TestChipAdaptorDigitalIO(t *testing.T) {
	a := initTestChipAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio408/value",
		"/sys/class/gpio/gpio408/direction",
		"/sys/class/gpio/gpio415/value",
		"/sys/class/gpio/gpio415/direction",
	})

	sysfs.SetFilesystem(fs)

	a.DigitalWrite("XIO-P0", 1)
	gobot.Assert(t, fs.Files["/sys/class/gpio/gpio408/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio415/value"].Contents = "1"
	i, _ := a.DigitalRead("XIO-P7")
	gobot.Assert(t, i, 1)

	gobot.Assert(t, a.DigitalWrite("XIO-P10", 1), errors.New("Not a valid pin"))
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
	gobot.Assert(t, data, []byte{0x00, 0x01})

	gobot.Assert(t, len(a.Finalize()), 0)
}
