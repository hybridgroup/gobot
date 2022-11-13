package joule

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
var _ sysfs.DigitalPinnerProvider = (*Adaptor)(nil)
var _ sysfs.PWMPinnerProvider = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

func initTestAdaptorWithMockedFilesystem() (*Adaptor, *sysfs.MockFilesystem) {
	a := NewAdaptor()
	mockPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm0/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm0/period",
		"/sys/class/pwm/pwmchip0/pwm0/enable",
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio13/value",
		"/sys/class/gpio/gpio13/direction",
		"/sys/class/gpio/gpio40/value",
		"/sys/class/gpio/gpio40/direction",
		"/sys/class/gpio/gpio446/value",
		"/sys/class/gpio/gpio446/direction",
		"/sys/class/gpio/gpio463/value",
		"/sys/class/gpio/gpio463/direction",
		"/sys/class/gpio/gpio421/value",
		"/sys/class/gpio/gpio421/direction",
		"/sys/class/gpio/gpio221/value",
		"/sys/class/gpio/gpio221/direction",
		"/sys/class/gpio/gpio243/value",
		"/sys/class/gpio/gpio243/direction",
		"/sys/class/gpio/gpio229/value",
		"/sys/class/gpio/gpio229/direction",
		"/sys/class/gpio/gpio253/value",
		"/sys/class/gpio/gpio253/direction",
		"/sys/class/gpio/gpio261/value",
		"/sys/class/gpio/gpio261/direction",
		"/sys/class/gpio/gpio214/value",
		"/sys/class/gpio/gpio214/direction",
		"/sys/class/gpio/gpio14/direction",
		"/sys/class/gpio/gpio14/value",
		"/sys/class/gpio/gpio165/direction",
		"/sys/class/gpio/gpio165/value",
		"/sys/class/gpio/gpio212/direction",
		"/sys/class/gpio/gpio212/value",
		"/sys/class/gpio/gpio213/direction",
		"/sys/class/gpio/gpio213/value",
		"/sys/class/gpio/gpio236/direction",
		"/sys/class/gpio/gpio236/value",
		"/sys/class/gpio/gpio237/direction",
		"/sys/class/gpio/gpio237/value",
		"/sys/class/gpio/gpio204/direction",
		"/sys/class/gpio/gpio204/value",
		"/sys/class/gpio/gpio205/direction",
		"/sys/class/gpio/gpio205/value",
		"/sys/class/gpio/gpio263/direction",
		"/sys/class/gpio/gpio263/value",
		"/sys/class/gpio/gpio262/direction",
		"/sys/class/gpio/gpio262/value",
		"/sys/class/gpio/gpio240/direction",
		"/sys/class/gpio/gpio240/value",
		"/sys/class/gpio/gpio241/direction",
		"/sys/class/gpio/gpio241/value",
		"/sys/class/gpio/gpio242/direction",
		"/sys/class/gpio/gpio242/value",
		"/sys/class/gpio/gpio218/direction",
		"/sys/class/gpio/gpio218/value",
		"/sys/class/gpio/gpio250/direction",
		"/sys/class/gpio/gpio250/value",
		"/sys/class/gpio/gpio451/direction",
		"/sys/class/gpio/gpio451/value",
		"/dev/i2c-0",
	}
	fs := a.sysfs.UseMockFilesystem(mockPaths)
	fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents = "5000"
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestName(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()

	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Joule"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestConnect(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()

	gobottest.Assert(t, a.GetDefaultBus(), 0)
}

func TestInvalidBus(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()
	gobottest.Assert(t, a.Connect(), nil)

	_, err := a.GetConnection(0xff, 10)
	gobottest.Assert(t, err, errors.New("Bus number 10 out of range"))
}

func TestFinalize(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()

	a.DigitalWrite("J12_1", 1)
	a.PwmWrite("J12_26", 100)

	gobottest.Assert(t, a.Finalize(), nil)

	_, err := a.GetConnection(0xff, 0)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestDigitalIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	a.DigitalWrite("J12_1", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio451/value"].Contents, "1")

	a.DigitalWrite("J12_1", 0)
	i, err := a.DigitalRead("J12_1")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i, 0)
}

func TestDigitalWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithWriteError = true

	err := a.DigitalWrite("13", 1)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestDigitalReadWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithWriteError = true

	_, err := a.DigitalRead("13")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestI2c(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()
	a.sysfs.UseMockSyscall()

	con, err := a.GetConnection(0xff, 0)
	gobottest.Assert(t, err, nil)

	_, err = con.Write([]byte{0x00, 0x01})
	gobottest.Assert(t, err, nil)

	data := []byte{42, 42}
	_, err = con.Read(data)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestPwm(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	err := a.PwmWrite("J12_26", 100)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "3921568")

	err = a.PwmWrite("4", 100)
	gobottest.Assert(t, err, errors.New("Not a valid pin"))

	err = a.PwmWrite("J12_1", 100)
	gobottest.Assert(t, err, errors.New("Not a PWM pin"))
}

func TestPwmPinExportError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	delete(fs.Files, "/sys/class/pwm/pwmchip0/export")

	err := a.PwmWrite("J12_26", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/pwm/pwmchip0/export: No such file"), true)
}

func TestPwmPinEnableError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	delete(fs.Files, "/sys/class/pwm/pwmchip0/pwm0/enable")

	err := a.PwmWrite("J12_26", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/pwm/pwmchip0/pwm0/enable: No such file"), true)
}
