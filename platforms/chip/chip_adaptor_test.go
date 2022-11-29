package chip

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/system"
)

// make sure that this Adaptor fulfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
var _ gobot.PWMPinnerProvider = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

var mockPaths = []string{
	"/sys/class/gpio/export",
	"/sys/class/gpio/unexport",
	"/sys/class/gpio/gpio50/value",
	"/sys/class/gpio/gpio50/direction",
	"/sys/class/gpio/gpio139/value",
	"/sys/class/gpio/gpio139/direction",
	"/sys/class/pwm/pwmchip0/export",
	"/sys/class/pwm/pwmchip0/unexport",
	"/sys/class/pwm/pwmchip0/pwm0/enable",
	"/sys/class/pwm/pwmchip0/pwm0/duty_cycle",
	"/sys/class/pwm/pwmchip0/pwm0/polarity",
	"/sys/class/pwm/pwmchip0/pwm0/period",
}

func initTestAdaptorWithMockedFilesystem() (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func initTestProAdaptorWithMockedFilesystem() (*Adaptor, *system.MockFilesystem) {
	a := NewProAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestName(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "CHIP"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestNewProAdaptor(t *testing.T) {
	a := NewProAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "CHIP Pro"), true)
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.DigitalWrite("CSID7", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.PwmWrite("PWM0", 100), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestDigitalIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	a.Connect()

	a.DigitalWrite("CSID7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio139/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio50/value"].Contents = "1"
	i, _ := a.DigitalRead("TWI2-SDA")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("XIO-P10", 1), errors.New("Not a valid pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestProDigitalIO(t *testing.T) {
	a, fs := initTestProAdaptorWithMockedFilesystem()
	a.Connect()

	a.DigitalWrite("CSID7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio139/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio50/value"].Contents = "1"
	i, _ := a.DigitalRead("TWI2-SDA")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("XIO-P0", 1), errors.New("Not a valid pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestI2c(t *testing.T) {
	a := NewAdaptor()
	a.sys.UseMockFilesystem([]string{"/dev/i2c-1"})
	a.sys.UseMockSyscall()

	a.Connect()

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

func TestInvalidPWMPin(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()
	a.Connect()

	err := a.PwmWrite("LCD-D2", 42)
	gobottest.Assert(t, err.Error(), "Not a PWM pin")

	err = a.ServoWrite("LCD-D2", 120)
	gobottest.Assert(t, err.Error(), "Not a PWM pin")
}

func TestPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	a.Connect()

	err := a.PwmWrite("PWM0", 100)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/export"].Contents, "0")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/enable"].Contents, "1")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "3921568")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/polarity"].Contents, "normal")

	err = a.ServoWrite("PWM0", 0)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "500000")

	err = a.ServoWrite("PWM0", 180)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "2000000")
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestPwmWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithWriteError = true

	err := a.PwmWrite("PWM0", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestPwmReadError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithReadError = true

	err := a.PwmWrite("PWM0", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "read error"), true)
}

func TestDefaultBus(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()
	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

func TestGetConnectionInvalidBus(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()
	_, err := a.GetConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Bus number 99 out of range"))
}
