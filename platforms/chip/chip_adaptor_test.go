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

// make sure that this Adaptor fullfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

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
		"/sys/class/gpio/gpio50/value",
		"/sys/class/gpio/gpio50/direction",
		"/sys/class/gpio/gpio139/value",
		"/sys/class/gpio/gpio139/direction",
	})

	sysfs.SetFilesystem(fs)

	a.DigitalWrite("CSID7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio139/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio50/value"].Contents = "1"
	i, _ := a.DigitalRead("TWI2-SDA")
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

	con, err := a.GetConnection(0xff, 1)
	gobottest.Assert(t, err, nil)

	con.Write([]byte{0x00, 0x01})
	data := []byte{42, 42}
	con.Read(data)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestChipAdaptorInvalidPWMPin(t *testing.T) {
	a := initTestChipAdaptor()

	err := a.PwmWrite("LCD-D2", 42)
	gobottest.Refute(t, err, nil)

	err = a.ServoWrite("LCD-D2", 120)
	gobottest.Refute(t, err, nil)
}

func TestChipAdaptorPWM(t *testing.T) {
	a := initTestChipAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm0/enable",
		"/sys/class/pwm/pwmchip0/pwm0/period",
		"/sys/class/pwm/pwmchip0/pwm0/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm0/polarity",
	})
	sysfs.SetFilesystem(fs)

	err := a.PwmWrite("PWM0", 123)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/export"].Contents, "0")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/enable"].Contents, "1")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents, "10000000")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "4823529")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/polarity"].Contents, "normal")

	err = a.ServoWrite("PWM0", 0)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "500000")

	err = a.ServoWrite("PWM0", 180)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "2000000")
}
