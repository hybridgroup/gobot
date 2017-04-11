package edison

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

// make sure that this Adaptor fullfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ aio.AnalogReader = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

func initTestAdaptor() (*Adaptor, *sysfs.MockFilesystem) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/bus/iio/devices/iio:device1/in_voltage0_raw",
		"/sys/kernel/debug/gpio_debug/gpio111/current_pinmux",
		"/sys/kernel/debug/gpio_debug/gpio115/current_pinmux",
		"/sys/kernel/debug/gpio_debug/gpio114/current_pinmux",
		"/sys/kernel/debug/gpio_debug/gpio109/current_pinmux",
		"/sys/kernel/debug/gpio_debug/gpio131/current_pinmux",
		"/sys/kernel/debug/gpio_debug/gpio129/current_pinmux",
		"/sys/kernel/debug/gpio_debug/gpio40/current_pinmux",
		"/sys/kernel/debug/gpio_debug/gpio13/current_pinmux",
		"/sys/kernel/debug/gpio_debug/gpio28/current_pinmux",
		"/sys/kernel/debug/gpio_debug/gpio27/current_pinmux",
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm1/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm1/period",
		"/sys/class/pwm/pwmchip0/pwm1/enable",
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio13/value",
		"/sys/class/gpio/gpio13/direction",
		"/sys/class/gpio/gpio40/value",
		"/sys/class/gpio/gpio40/direction",
		"/sys/class/gpio/gpio128/value",
		"/sys/class/gpio/gpio128/direction",
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
		"/dev/i2c-6",
	})
	sysfs.SetFilesystem(fs)
	fs.Files["/sys/class/pwm/pwmchip0/pwm1/period"].Contents = "5000\n"
	a.Connect()
	return a, fs
}

func TestEdisonAdaptorName(t *testing.T) {
	a, _ := initTestAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Edison"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestAdaptorConnect(t *testing.T) {
	a, _ := initTestAdaptor()
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.GetDefaultBus(), 6)

	a = NewAdaptor()
	sysfs.SetFilesystem(sysfs.NewMockFilesystem([]string{}))
	gobottest.Refute(t, a.Connect(), nil)
}

func TestAdaptorConnectSparkfun(t *testing.T) {
	a, _ := initTestAdaptor()
	a.SetBoard("sparkfun")
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

func TestAdaptorConnectMiniboard(t *testing.T) {
	a, _ := initTestAdaptor()
	a.SetBoard("miniboard")
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

func TestAdaptorConnectUnknown(t *testing.T) {
	a, _ := initTestAdaptor()
	a.SetBoard("wha")
	gobottest.Refute(t, a.Connect(), nil)
}

func TestAdaptorFinalize(t *testing.T) {
	a, _ := initTestAdaptor()
	a.DigitalWrite("3", 1)
	a.PwmWrite("5", 100)

	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.GetConnection(0xff, 6)

	gobottest.Assert(t, a.Finalize(), nil)

	sysfs.SetFilesystem(sysfs.NewMockFilesystem([]string{}))
	gobottest.Refute(t, a.Finalize(), nil)
}

func TestAdaptorDigitalIO(t *testing.T) {
	a, fs := initTestAdaptor()

	a.DigitalWrite("13", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio40/value"].Contents, "1")

	a.DigitalWrite("2", 0)
	i, err := a.DigitalRead("2")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i, 0)
}

func TestAdaptorI2c(t *testing.T) {
	a, _ := initTestAdaptor()

	sysfs.SetSyscall(&sysfs.MockSyscall{})
	con, err := a.GetConnection(0xff, 6)
	gobottest.Assert(t, err, nil)

	con.Write([]byte{0x00, 0x01})
	data := []byte{42, 42}
	con.Read(data)
	gobottest.Assert(t, data, []byte{0x00, 0x01})
}

func TestAdaptorPwm(t *testing.T) {
	a, fs := initTestAdaptor()

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm1/duty_cycle"].Contents, "1960")

	err = a.PwmWrite("7", 100)
	gobottest.Assert(t, err, errors.New("Not a PWM pin"))
}

func TestAdaptorAnalog(t *testing.T) {
	a, fs := initTestAdaptor()

	fs.Files["/sys/bus/iio/devices/iio:device1/in_voltage0_raw"].Contents = "1000\n"
	i, _ := a.AnalogRead("0")
	gobottest.Assert(t, i, 250)
}
