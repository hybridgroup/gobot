package up2

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/spi"
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
var _ spi.Connector = (*Adaptor)(nil)

func initTestUP2Adaptor() (*Adaptor, *sysfs.MockFilesystem) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio462/value",
		"/sys/class/gpio/gpio462/direction",
		"/sys/class/gpio/gpio432/value",
		"/sys/class/gpio/gpio432/direction",
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm0/enable",
		"/sys/class/pwm/pwmchip0/pwm0/period",
		"/sys/class/pwm/pwmchip0/pwm0/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm0/polarity",
		"/sys/class/leds/upboard:green:/brightness",
	})

	sysfs.SetFilesystem(fs)
	return a, fs
}

func TestUP2AdaptorName(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "UP2"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestUP2AdaptorDigitalIO(t *testing.T) {
	a, fs := initTestUP2Adaptor()
	a.Connect()

	a.DigitalWrite("7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio462/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio432/value"].Contents = "1"
	i, _ := a.DigitalRead("13")
	gobottest.Assert(t, i, 1)

	a.DigitalWrite("green", 1)
	gobottest.Assert(t,
		fs.Files["/sys/class/leds/upboard:green:/brightness"].Contents,
		"1",
	)

	gobottest.Assert(t, a.DigitalWrite("99", 1), errors.New("Not a valid pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestAdaptorDigitalWriteError(t *testing.T) {
	a, fs := initTestUP2Adaptor()
	fs.WithWriteError = true

	err := a.DigitalWrite("7", 1)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestAdaptorDigitalReadWriteError(t *testing.T) {
	a, fs := initTestUP2Adaptor()
	fs.WithWriteError = true

	_, err := a.DigitalRead("7")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestUP2AdaptorI2c(t *testing.T) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-5",
	})
	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})

	con, err := a.GetConnection(0xff, 5)
	gobottest.Assert(t, err, nil)

	con.Write([]byte{0x00, 0x01})
	data := []byte{42, 42}
	con.Read(data)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestAdaptorSPI(t *testing.T) {
	a := NewAdaptor()
	a.Connect()

	gobottest.Assert(t, a.GetSpiDefaultBus(), 0)
	gobottest.Assert(t, a.GetSpiDefaultMode(), 0)
	gobottest.Assert(t, a.GetSpiDefaultMaxSpeed(), int64(500000))

	_, err := a.GetSpiConnection(10, 0, 0, 8, 500000)
	gobottest.Assert(t, err.Error(), "Bus number 10 out of range")

	// TODO: test tx/rx here...

}

func TestUP2AdaptorInvalidPWMPin(t *testing.T) {
	a, _ := initTestUP2Adaptor()
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

func TestUP2AdaptorPWM(t *testing.T) {
	a, fs := initTestUP2Adaptor()

	err := a.PwmWrite("32", 100)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/export"].Contents, "0")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/enable"].Contents, "1")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "3921568")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/polarity"].Contents, "normal")

	err = a.ServoWrite("32", 0)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "500000")

	err = a.ServoWrite("32", 180)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "2000000")
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestUP2AdaptorPwmWriteError(t *testing.T) {
	a, fs := initTestUP2Adaptor()
	fs.WithWriteError = true

	err := a.PwmWrite("32", 100)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestUP2AdaptorPwmReadError(t *testing.T) {
	a, fs := initTestUP2Adaptor()
	fs.WithReadError = true

	err := a.PwmWrite("32", 100)
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestUP2I2CDefaultBus(t *testing.T) {
	a, _ := initTestUP2Adaptor()
	gobottest.Assert(t, a.GetDefaultBus(), 5)
}

func TestUP2GetConnectionInvalidBus(t *testing.T) {
	a, _ := initTestUP2Adaptor()
	_, err := a.GetConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Bus number 99 out of range"))
}

func TestUP2FinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestUP2Adaptor()

	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.DigitalWrite("7", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestUP2FinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestUP2Adaptor()

	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.PwmWrite("32", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}
