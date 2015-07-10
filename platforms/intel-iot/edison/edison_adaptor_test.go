package edison

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

type NullReadWriteCloser struct {
	contents []byte
}

func (n *NullReadWriteCloser) SetAddress(int) error {
	return nil
}

func (n *NullReadWriteCloser) Write(b []byte) (int, error) {
	n.contents = make([]byte, len(b))
	copy(n.contents[:], b[:])

	return len(b), nil
}

func (n *NullReadWriteCloser) Read(b []byte) (int, error) {
	copy(b, n.contents)
	return len(b), nil
}

var closeErr error = nil

func (n *NullReadWriteCloser) Close() error {
	return closeErr
}

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

func TestEdisonAdaptor(t *testing.T) {
	a, _ := initTestEdisonAdaptor()
	gobot.Assert(t, a.Name(), "myAdaptor")
}
func TestEdisonAdaptorConnect(t *testing.T) {
	a, _ := initTestEdisonAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)

	a = NewEdisonAdaptor("myAdaptor")
	sysfs.SetFilesystem(sysfs.NewMockFilesystem([]string{}))
	gobot.Refute(t, len(a.Connect()), 0)
}

func TestEdisonAdaptorFinalize(t *testing.T) {
	a, _ := initTestEdisonAdaptor()
	a.DigitalWrite("3", 1)
	a.PwmWrite("5", 100)

	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.I2cStart(0xff)

	gobot.Assert(t, len(a.Finalize()), 0)

	closeErr = errors.New("close error")
	sysfs.SetFilesystem(sysfs.NewMockFilesystem([]string{}))
	gobot.Refute(t, len(a.Finalize()), 0)
}

func TestEdisonAdaptorDigitalIO(t *testing.T) {
	a, fs := initTestEdisonAdaptor()

	a.DigitalWrite("13", 1)
	gobot.Assert(t, fs.Files["/sys/class/gpio/gpio40/value"].Contents, "1")

	a.DigitalWrite("2", 0)
	i, err := a.DigitalRead("2")
	gobot.Assert(t, err, nil)
	gobot.Assert(t, i, 0)
}

func TestEdisonAdaptorI2c(t *testing.T) {
	a, _ := initTestEdisonAdaptor()

	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.I2cStart(0xff)

	a.i2cDevice = &NullReadWriteCloser{}
	a.I2cWrite(0xff, []byte{0x00, 0x01})

	data, _ := a.I2cRead(0xff, 2)
	gobot.Assert(t, data, []byte{0x00, 0x01})
}

func TestEdisonAdaptorPwm(t *testing.T) {
	a, fs := initTestEdisonAdaptor()

	err := a.PwmWrite("5", 100)
	gobot.Assert(t, err, nil)
	gobot.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm1/duty_cycle"].Contents, "1960")

	err = a.PwmWrite("7", 100)
	gobot.Assert(t, err, errors.New("Not a PWM pin"))
}

func TestEdisonAdaptorAnalog(t *testing.T) {
	a, fs := initTestEdisonAdaptor()

	fs.Files["/sys/bus/iio/devices/iio:device1/in_voltage0_raw"].Contents = "1000\n"
	i, _ := a.AnalogRead("0")
	gobot.Assert(t, i, 250)
}
