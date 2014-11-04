package beaglebone

import (
	"io"
	"os"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

func initTestBeagleboneAdaptor() *BeagleboneAdaptor {
	i2cLocation = os.DevNull
	sysfs.WriteFile = func(path string, data []byte) (i int, err error) {
		return
	}
	a := NewBeagleboneAdaptor("myAdaptor")
	a.connect = func() {}
	a.Connect()
	a.DigitalWrite("P9_12", 1)
	a.i2cDevice = new(gobot.NullReadWriteCloser)
	return a
}

func TestBeagleboneAdaptorFinalize(t *testing.T) {
	gobot.Assert(t, initTestBeagleboneAdaptor().Finalize(), true)
}

func TestBeagleboneAdaptorDigitalIO(t *testing.T) {
	a := initTestBeagleboneAdaptor()
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

	a.DigitalWrite("P9_12", 1)
	gobot.Assert(t, lastWritePath, "/sys/class/gpio/gpio60/value")
	gobot.Assert(t, lastWriteData, []byte{49})

	i := a.DigitalRead("P8_31")
	gobot.Assert(t, lastReadPath, "/sys/class/gpio/gpio10/value")
	gobot.Assert(t, i, 1)
}

func TestBeagleboneAdaptorI2c(t *testing.T) {
	a := initTestBeagleboneAdaptor()
	a.I2cStart(0xff)
	var _ io.ReadWriteCloser = a.i2cDevice

	a.i2cDevice = new(gobot.NullReadWriteCloser)
	a.I2cWrite([]byte{0x00, 0x01})
	gobot.Assert(t, a.I2cRead(2), make([]byte, 2))
}
