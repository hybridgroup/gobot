package chip

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

func initTestChipAdaptor() (*Adaptor, *sysfs.MockFilesystem) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
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
	})

	sysfs.SetFilesystem(fs)
	return a, fs
}

func initTestChipProAdaptor() (*Adaptor, *sysfs.MockFilesystem) {
	a := NewProAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
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
	})

	sysfs.SetFilesystem(fs)
	return a, fs
}

func TestChipAdaptorName(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "CHIP"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestChipAdaptorBoard(t *testing.T) {
	a := NewAdaptor()
	a.SetBoard("pro")
	gobottest.Assert(t, a.board, "pro")

	gobottest.Assert(t, a.SetBoard("bad"), errors.New("Invalid board type"))
}

func TestAdaptorFinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestChipAdaptor()
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.DigitalWrite("CSID7", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestAdaptorFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestChipAdaptor()
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.PwmWrite("PWM0", 100), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestChipAdaptorDigitalIO(t *testing.T) {
	a, fs := initTestChipAdaptor()
	a.Connect()

	a.DigitalWrite("CSID7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio139/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio50/value"].Contents = "1"
	i, _ := a.DigitalRead("TWI2-SDA")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("XIO-P10", 1), errors.New("Not a valid pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestChipProAdaptorDigitalIO(t *testing.T) {
	a, fs := initTestChipProAdaptor()
	a.Connect()

	a.DigitalWrite("CSID7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio139/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio50/value"].Contents = "1"
	i, _ := a.DigitalRead("TWI2-SDA")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("XIO-P0", 1), errors.New("Not a valid pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestAdaptorDigitalWriteError(t *testing.T) {
	a, fs := initTestChipAdaptor()
	fs.WithWriteError = true

	err := a.DigitalWrite("CSID7", 1)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestAdaptorDigitalReadWriteError(t *testing.T) {
	a, fs := initTestChipAdaptor()
	fs.WithWriteError = true

	_, err := a.DigitalRead("CSID7")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestChipAdaptorI2c(t *testing.T) {
	a := NewAdaptor()
	a.Connect()

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

func TestChipAdaptorInvalidPWMPin(t *testing.T) {
	a, _ := initTestChipAdaptor()
	a.Connect()

	err := a.PwmWrite("LCD-D2", 42)
	gobottest.Refute(t, err, nil)

	err = a.ServoWrite("LCD-D2", 120)
	gobottest.Refute(t, err, nil)
}

func TestChipAdaptorPWM(t *testing.T) {
	a, fs := initTestChipAdaptor()
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

func TestAdaptorPwmWriteError(t *testing.T) {
	a, fs := initTestChipAdaptor()
	fs.WithWriteError = true

	err := a.PwmWrite("PWM0", 100)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestAdaptorPwmReadError(t *testing.T) {
	a, fs := initTestChipAdaptor()
	fs.WithReadError = true

	err := a.PwmWrite("PWM0", 100)
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestChipDefaultBus(t *testing.T) {
	a, _ := initTestChipAdaptor()
	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

func TestChipGetConnectionInvalidBus(t *testing.T) {
	a, _ := initTestChipAdaptor()
	_, err := a.GetConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Bus number 99 out of range"))
}
