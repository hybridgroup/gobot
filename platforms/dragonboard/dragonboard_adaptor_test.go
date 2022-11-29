package dragonboard

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
)

// make sure that this Adaptor fulfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

func initTestAdaptor(t *testing.T) *Adaptor {
	a := NewAdaptor()
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a
}

func TestName(t *testing.T) {
	a := initTestAdaptor(t)
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "DragonBoard"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestDigitalIO(t *testing.T) {
	a := initTestAdaptor(t)
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio36/value",
		"/sys/class/gpio/gpio36/direction",
		"/sys/class/gpio/gpio12/value",
		"/sys/class/gpio/gpio12/direction",
	}
	fs := a.sys.UseMockFilesystem(mockPaths)

	_ = a.DigitalWrite("GPIO_B", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio12/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio36/value"].Contents = "1"
	i, _ := a.DigitalRead("GPIO_A")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("GPIO_M", 1), errors.New("Not a valid pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestI2c(t *testing.T) {
	a := initTestAdaptor(t)
	a.sys.UseMockFilesystem([]string{"/dev/i2c-1"})
	a.sys.UseMockSyscall()

	con, err := a.GetConnection(0xff, 1)
	gobottest.Assert(t, err, nil)

	_, err = con.Write([]byte{0x00, 0x01})
	gobottest.Assert(t, err, nil)

	data := []byte{42, 42}
	_, err = con.Read(data)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestDefaultBus(t *testing.T) {
	a := initTestAdaptor(t)
	gobottest.Assert(t, a.GetDefaultBus(), 0)
}

func TestGetConnectionInvalidBus(t *testing.T) {
	a := initTestAdaptor(t)
	_, err := a.GetConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Bus number 99 out of range"))
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a := initTestAdaptor(t)
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio36/value",
		"/sys/class/gpio/gpio36/direction",
		"/sys/class/gpio/gpio12/value",
		"/sys/class/gpio/gpio12/direction",
	}
	fs := a.sys.UseMockFilesystem(mockPaths)

	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.DigitalWrite("GPIO_B", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}
