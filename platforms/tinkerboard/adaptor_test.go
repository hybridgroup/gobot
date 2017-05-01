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
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ sysfs.DigitalPinnerProvider = (*Adaptor)(nil)
var _ sysfs.PWMPinnerProvider = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

func initTestTinkerboardAdaptor() (*Adaptor, *sysfs.MockFilesystem) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio17/value",
		"/sys/class/gpio/gpio17/direction",
		"/sys/class/gpio/gpio160/value",
		"/sys/class/gpio/gpio160/direction",
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm0/enable",
		"/sys/class/pwm/pwmchip0/pwm0/period",
		"/sys/class/pwm/pwmchip0/pwm0/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm0/polarity",
	})

	sysfs.SetFilesystem(fs)
	return a, fs
}

func TestTinkerboardAdaptorName(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Tinker Board"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestTinkerboardAdaptorDigitalIO(t *testing.T) {
	a, fs := initTestTinkerboardAdaptor()
	a.Connect()

	a.DigitalWrite("7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio17/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio160/value"].Contents = "1"
	i, _ := a.DigitalRead("10")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("99", 1), errors.New("Not a valid pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestAdaptorDigitalWriteError(t *testing.T) {
	a, fs := initTestTinkerboardAdaptor()
	fs.WithWriteError = true

	err := a.DigitalWrite("7", 1)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestAdaptorDigitalReadWriteError(t *testing.T) {
	a, fs := initTestTinkerboardAdaptor()
	fs.WithWriteError = true

	_, err := a.DigitalRead("7")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestTinkerboardAdaptorI2c(t *testing.T) {
	a := NewAdaptor()
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

func TestTinkerboardAdaptorInvalidPWMPin(t *testing.T) {
	a, _ := initTestTinkerboardAdaptor()
	a.Connect()

	err := a.PwmWrite("666", 42)
	gobottest.Refute(t, err, nil)

	err = a.ServoWrite("666", 120)
	gobottest.Refute(t, err, nil)

	err = a.PwmWrite("3", 42)
	gobottest.Refute(t, err, nil)

	err = a.ServoWrite("3", 120)
	gobottest.Refute(t, err, nil)
}

func TestTinkerboardAdaptorPWM(t *testing.T) {
	a, fs := initTestTinkerboardAdaptor()

	err := a.PwmWrite("33", 100)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/export"].Contents, "0")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/enable"].Contents, "1")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "3921568")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/polarity"].Contents, "normal")

	err = a.ServoWrite("33", 0)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "500000")

	err = a.ServoWrite("33", 180)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "2000000")
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestTinkerboardAdaptorPwmWriteError(t *testing.T) {
	a, fs := initTestTinkerboardAdaptor()
	fs.WithWriteError = true

	err := a.PwmWrite("33", 100)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestTinkerboardAdaptorPwmReadError(t *testing.T) {
	a, fs := initTestTinkerboardAdaptor()
	fs.WithReadError = true

	err := a.PwmWrite("33", 100)
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestTinkerboardDefaultBus(t *testing.T) {
	a, _ := initTestTinkerboardAdaptor()
	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

func TestTinkerboardGetConnectionInvalidBus(t *testing.T) {
	a, _ := initTestTinkerboardAdaptor()
	_, err := a.GetConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Bus number 99 out of range"))
}

func TestTinkerboardFinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestTinkerboardAdaptor()

	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.DigitalWrite("7", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestTinkerboardFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestTinkerboardAdaptor()

	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.PwmWrite("33", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}
