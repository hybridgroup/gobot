package joule

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
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

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

var closeErr error

func (n *NullReadWriteCloser) Close() error {
	return closeErr
}

func initTestAdaptor() (*Adaptor, *sysfs.MockFilesystem) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
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

func TestAdaptorConnect(t *testing.T) {
	a, _ := initTestAdaptor()
	gobottest.Assert(t, a.Connect(), nil)
}

func TestAdaptorFinalize(t *testing.T) {
	a, _ := initTestAdaptor()
	a.DigitalWrite("1", 1)
	a.PwmWrite("25", 100)

	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.GetConnection(0xff, 0)

	gobottest.Assert(t, a.Finalize(), nil)

	closeErr = errors.New("close error")
	sysfs.SetFilesystem(sysfs.NewMockFilesystem([]string{}))
	gobottest.Refute(t, a.Finalize(), nil)
}

func TestAdaptorDigitalIO(t *testing.T) {
	a, fs := initTestAdaptor()

	a.DigitalWrite("1", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio446/value"].Contents, "1")

	a.DigitalWrite("2", 0)
	i, err := a.DigitalRead("2")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i, 0)
}

func TestAdaptorI2c(t *testing.T) {
	a, _ := initTestAdaptor()

	sysfs.SetSyscall(&sysfs.MockSyscall{})
	con, err := a.GetConnection(0xff, 0)
	gobottest.Assert(t, err, nil)

	con.Write([]byte{0x00, 0x01})
	data := []byte{42, 42}
	con.Read(data)
	gobottest.Assert(t, data, []byte{0x00, 0x01})
}

func TestAdaptorPwm(t *testing.T) {
	a, fs := initTestAdaptor()

	err := a.PwmWrite("25", 100)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "1960")

	err = a.PwmWrite("4", 100)
	gobottest.Assert(t, err, errors.New("Not a PWM pin"))
}
