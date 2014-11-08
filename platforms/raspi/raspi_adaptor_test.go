package raspi

import (
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

func initTestRaspiAdaptor() *RaspiAdaptor {
	boardRevision = func() (string, string) {
		return "3", "/dev/i2c-1"
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
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio4/value",
		"/sys/class/gpio/gpio4/direction",
		"/sys/class/gpio/gpio27/value",
		"/sys/class/gpio/gpio27/direction",
	})

	sysfs.SetFilesystem(fs)

	a.DigitalWrite("7", 1)
	gobot.Assert(t, fs.Files["/sys/class/gpio/gpio4/value"].Contents, "1")

	a.DigitalWrite("13", 1)
	i := a.DigitalRead("13")
	gobot.Assert(t, i, 1)
}

func TestRaspiAdaptorI2c(t *testing.T) {
	a := initTestRaspiAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.I2cStart(0xff)

	a.I2cWrite([]byte{0x00, 0x01})
	gobot.Assert(t, a.I2cRead(2), []byte{0x00, 0x01})
}
