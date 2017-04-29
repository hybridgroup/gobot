package dragonboard

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

func initTestDragonBoardAdaptor(t *testing.T) *Adaptor {
	a := NewAdaptor()
	if err := a.Connect(); err != nil {
		t.Error(err)
	}
	return a
}

func TestDragonBoardAdaptorName(t *testing.T) {
	a := initTestDragonBoardAdaptor(t)
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "DragonBoard"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestDragonBoardAdaptorDigitalIO(t *testing.T) {
	a := initTestDragonBoardAdaptor(t)
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio36/value",
		"/sys/class/gpio/gpio36/direction",
		"/sys/class/gpio/gpio12/value",
		"/sys/class/gpio/gpio12/direction",
	})

	sysfs.SetFilesystem(fs)

	_ = a.DigitalWrite("GPIO_B", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio12/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio36/value"].Contents = "1"
	i, _ := a.DigitalRead("GPIO_A")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("GPIO_M", 1), errors.New("Not a valid pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestDragonBoardAdaptorI2c(t *testing.T) {
	a := initTestDragonBoardAdaptor(t)
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})

	con, err := a.GetConnection(0xff, 1)
	gobottest.Assert(t, err, nil)

	_, _ = con.Write([]byte{0x00, 0x01})
	data := []byte{42, 42}
	_, _ = con.Read(data)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestDragonBoardDefaultBus(t *testing.T) {
	a := initTestDragonBoardAdaptor(t)
	gobottest.Assert(t, a.GetDefaultBus(), 0)
}

func TestDragonBoardGetConnectionInvalidBus(t *testing.T) {
	a := initTestDragonBoardAdaptor(t)
	_, err := a.GetConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Bus number 99 out of range"))
}

func TestAdaptorFinalizeErrorAfterGPIO(t *testing.T) {
	a := initTestDragonBoardAdaptor(t)
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio36/value",
		"/sys/class/gpio/gpio36/direction",
		"/sys/class/gpio/gpio12/value",
		"/sys/class/gpio/gpio12/direction",
	})

	sysfs.SetFilesystem(fs)

	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.DigitalWrite("GPIO_B", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}
