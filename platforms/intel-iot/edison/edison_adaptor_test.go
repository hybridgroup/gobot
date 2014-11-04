package edison

import (
	"io"
	"os"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

func initTestEdisonAdaptor() *EdisonAdaptor {
	i2cLocation = os.DevNull
	sysfs.WriteFile = func(path string, data []byte) (i int, err error) {
		return
	}
	writeFile = func(name, data string) error {
		return nil
	}
	readFile = func(name string) ([]byte, error) {
		return []byte("11"), nil
	}
	a := NewEdisonAdaptor("myAdaptor")
	a.Connect()
	return a
}

func TestEdisonAdaptorFinalize(t *testing.T) {
	a := initTestEdisonAdaptor()
	a.DigitalWrite("3", 1)
	a.PwmWrite("5", 100)
	a.i2cDevice = new(gobot.NullReadWriteCloser)
	gobot.Assert(t, a.Finalize(), true)
}

func TestEdisonAdaptorDigitalIO(t *testing.T) {
	a := initTestEdisonAdaptor()
	lastWritePath := ""
	lastReadPath := ""
	lastWriteData := []byte{}

	sysfs.WriteFile = func(path string, data []byte) (i int, err error) {
		lastWritePath = path
		lastWriteData = data
		return
	}
	sysfs.ReadFile = func(path string) (b []byte, err error) {
		lastReadPath = path
		return []byte("1"), nil
	}

	a.DigitalWrite("13", 1)
	gobot.Assert(t, lastWritePath, "/sys/class/gpio/gpio40/value")
	gobot.Assert(t, lastWriteData, []byte{49})

	i := a.DigitalRead("2")
	gobot.Assert(t, lastReadPath, "/sys/class/gpio/gpio128/value")
	gobot.Assert(t, i, 1)
}

func TestEdisonAdaptorI2c(t *testing.T) {
	a := initTestEdisonAdaptor()
	a.I2cStart(0xff)
	var _ io.ReadWriteCloser = a.i2cDevice

	a.i2cDevice = new(gobot.NullReadWriteCloser)
	a.I2cWrite([]byte{0x00, 0x01})
	gobot.Assert(t, a.I2cRead(2), make([]byte, 2))
}

func TestEdisonAdaptorPwm(t *testing.T) {
	a := initTestEdisonAdaptor()
	lastWritePath := ""
	lastReadPath := ""
	lastWriteData := ""

	writeFile = func(path, data string) (err error) {
		lastWritePath = path
		lastWriteData = data
		return
	}

	readFile = func(path string) (b []byte, err error) {
		lastReadPath = path
		return []byte("100\n"), nil
	}
	a.PwmWrite("5", 100)
	gobot.Assert(t, lastWritePath, "/sys/class/pwm/pwmchip0/pwm1/duty_cycle")
}

func TestEdisonAdaptorAnalog(t *testing.T) {
	a := initTestEdisonAdaptor()
	lastReadPath := ""
	readFile = func(path string) (b []byte, err error) {
		lastReadPath = path
		return []byte("100\n"), nil
	}
	i := a.AnalogRead("0")
	gobot.Assert(t, lastReadPath, "/sys/bus/iio/devices/iio:device1/in_voltage0_raw")
	gobot.Assert(t, i, 100)
}
