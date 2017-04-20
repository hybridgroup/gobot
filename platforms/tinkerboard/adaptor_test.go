package tinkerboard

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

// make sure that this Adaptor fullfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

func initTestTinkerboardAdaptor() *Adaptor {
	a := NewAdaptor()
	a.Connect()
	return a
}

func TestTinkerboardAdaptorName(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Tinkerboard"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestTinkerboardAdaptorDigitalIO(t *testing.T) {
	a := initTestTinkerboardAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio17/value",
		"/sys/class/gpio/gpio17/direction",
		"/sys/class/gpio/gpio160/value",
		"/sys/class/gpio/gpio160/direction",
	})

	sysfs.SetFilesystem(fs)

	a.DigitalWrite("7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio17/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio160/value"].Contents = "1"
	i, _ := a.DigitalRead("10")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("99", 1), errors.New("Not a valid pin"))
}

func TestTinkerboardAdaptorI2c(t *testing.T) {
	a := initTestTinkerboardAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})

	con, err := a.GetConnection(0xff, 1)
	gobottest.Assert(t, err, nil)

	con.Write([]byte{0x00, 0x01})
	data := []byte{42, 42}
	con.Read(data)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestTinkerboardDefaultBus(t *testing.T) {
	a := initTestTinkerboardAdaptor()
	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

func TestTinkerboardGetConnectionInvalidBus(t *testing.T) {
	a := initTestTinkerboardAdaptor()
	_, err := a.GetConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Bus number 99 out of range"))
}
