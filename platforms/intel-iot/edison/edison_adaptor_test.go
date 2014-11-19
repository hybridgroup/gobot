package edison

import (
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

func initTestEdisonAdaptor() (*EdisonAdaptor, *sysfs.MockFilesystem) {
	a := NewEdisonAdaptor("myAdaptor")
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
		"/sys/class/gpio/gpio229/direction",
		"/sys/class/gpio/gpio253/value",
		"/sys/class/gpio/gpio253/direction",
		"/sys/class/gpio/gpio261/value",
		"/sys/class/gpio/gpio261/direction",
		"/dev/i2c-6",
	})
	sysfs.SetFilesystem(fs)
	fs.Files["/sys/class/pwm/pwmchip0/pwm1/period"].Contents = "5000\n"
	a.Connect()
	return a, fs
}

func TestEdisonAdaptorFinalize(t *testing.T) {
	a, _ := initTestEdisonAdaptor()
	a.DigitalWrite("3", 1)
	a.PwmWrite("5", 100)
	a.i2cDevice = new(gobot.NullReadWriteCloser)
	gobot.Assert(t, a.Finalize(), nil)
}

func TestEdisonAdaptorDigitalIO(t *testing.T) {
	a, fs := initTestEdisonAdaptor()

	a.DigitalWrite("13", 1)
	gobot.Assert(t, fs.Files["/sys/class/gpio/gpio40/value"].Contents, "1")

	a.DigitalWrite("2", 0)
	i, _ := a.DigitalRead("2")
	gobot.Assert(t, i, 0)
}

func TestEdisonAdaptorI2c(t *testing.T) {
	a, _ := initTestEdisonAdaptor()

	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.I2cStart(0xff)

	a.I2cWrite([]byte{0x00, 0x01})
	gobot.Assert(t, a.I2cRead(2), []byte{0x00, 0x01})
}

func TestEdisonAdaptorPwm(t *testing.T) {
	a, fs := initTestEdisonAdaptor()

	a.PwmWrite("5", 100)
	gobot.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm1/duty_cycle"].Contents, "1960")
}

func TestEdisonAdaptorAnalog(t *testing.T) {
	a, fs := initTestEdisonAdaptor()

	fs.Files["/sys/bus/iio/devices/iio:device1/in_voltage0_raw"].Contents = "1000\n"
	i, _ := a.AnalogRead("0")
	gobot.Assert(t, i, 1000)
}
