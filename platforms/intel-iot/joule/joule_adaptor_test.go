package joule

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/sysfs"
)

var _ gobot.Adaptor = (*JouleAdaptor)(nil)

var _ gpio.DigitalReader = (*JouleAdaptor)(nil)
var _ gpio.DigitalWriter = (*JouleAdaptor)(nil)
var _ gpio.PwmWriter = (*JouleAdaptor)(nil)

var _ i2c.I2c = (*JouleAdaptor)(nil)

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

func initTestJouleAdaptor() (*JouleAdaptor, *sysfs.MockFilesystem) {
	a := NewJouleAdaptor("myAdaptor")
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
		"/dev/i2c-0",
	})
	sysfs.SetFilesystem(fs)
	fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents = "5000\n"
	a.Connect()
	return a, fs
}

func TestJouleAdaptor(t *testing.T) {
	a, _ := initTestJouleAdaptor()
	gobottest.Assert(t, a.Name(), "myAdaptor")
}

func TestJouleAdaptorConnect(t *testing.T) {
	a, _ := initTestJouleAdaptor()
	gobottest.Assert(t, len(a.Connect()), 0)
}

func TestJouleAdaptorFinalize(t *testing.T) {
	a, _ := initTestJouleAdaptor()
	a.DigitalWrite("1", 1)
	a.PwmWrite("25", 100)

	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.I2cStart(0xff)

	gobottest.Assert(t, len(a.Finalize()), 0)

	closeErr = errors.New("close error")
	sysfs.SetFilesystem(sysfs.NewMockFilesystem([]string{}))
	gobottest.Refute(t, len(a.Finalize()), 0)
}

func TestJouleAdaptorDigitalIO(t *testing.T) {
	a, fs := initTestJouleAdaptor()

	a.DigitalWrite("1", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio446/value"].Contents, "1")

	a.DigitalWrite("2", 0)
	i, err := a.DigitalRead("2")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i, 0)
}

func TestJouleAdaptorI2c(t *testing.T) {
	a, _ := initTestJouleAdaptor()

	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.I2cStart(0xff)

	a.i2cDevice = &NullReadWriteCloser{}
	a.I2cWrite(0xff, []byte{0x00, 0x01})

	data, _ := a.I2cRead(0xff, 2)
	gobottest.Assert(t, data, []byte{0x00, 0x01})
}

func TestJouleAdaptorPwm(t *testing.T) {
	a, fs := initTestJouleAdaptor()

	err := a.PwmWrite("25", 100)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "1960")

	err = a.PwmWrite("4", 100)
	gobottest.Assert(t, err, errors.New("Not a PWM pin"))
}
