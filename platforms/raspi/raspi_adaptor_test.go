package raspi

import (
	"io"
	"os"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

func initTestRaspiAdaptor() *RaspiAdaptor {
	i2cLocationFor = func(rev string) string {
		return os.DevNull
	}
	boardRevision = func() string {
		return "3"
	}
	sysfs.WriteFile = func(path string, data []byte) (i int, err error) {
		return
	}
	a := NewRaspiAdaptor("myAdaptor")
	a.Connect()
	return a
}

func TestRaspiAdaptorFinalize(t *testing.T) {
	a := initTestRaspiAdaptor()
	a.DigitalWrite("3", 1)
	a.i2cDevice = new(gobot.NullReadWriteCloser)
	gobot.Assert(t, a.Finalize(), true)
}

func TestRaspiAdaptorDigitalIO(t *testing.T) {
	a := initTestRaspiAdaptor()
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

	a.DigitalWrite("7", 1)
	gobot.Assert(t, lastWritePath, "/sys/class/gpio/gpio4/value")
	gobot.Assert(t, lastWriteData, []byte{49})

	i := a.DigitalRead("13")
	gobot.Assert(t, lastReadPath, "/sys/class/gpio/gpio27/value")
	gobot.Assert(t, i, 1)
}

func TestRaspiAdaptorI2c(t *testing.T) {
	a := initTestRaspiAdaptor()
	a.I2cStart(0xff)
	var _ io.ReadWriteCloser = a.i2cDevice

	a.i2cDevice = new(gobot.NullReadWriteCloser)
	a.I2cWrite([]byte{0x00, 0x01})
	gobot.Assert(t, a.I2cRead(2), make([]byte, 2))
}
