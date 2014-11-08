package beaglebone

import (
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

func initTestBeagleboneAdaptor() *BeagleboneAdaptor {
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

	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio60/value",
		"/sys/class/gpio/gpio60/direction",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
	})

	sysfs.SetFilesystem(fs)
	a.DigitalWrite("P9_12", 1)
	gobot.Assert(t, fs.Files["/sys/class/gpio/gpio60/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio10/value"].Contents = "1"
	i := a.DigitalRead("P8_31")
	gobot.Assert(t, i, 1)
}

func TestBeagleboneAdaptorI2c(t *testing.T) {
	a := initTestBeagleboneAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.I2cStart(0xff)

	a.I2cWrite([]byte{0x00, 0x01})
	gobot.Assert(t, a.I2cRead(2), []byte{0x00, 0x01})
}
