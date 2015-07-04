package raspi

import (
	"strings"
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

func initTestRaspiAdaptor() *RaspiAdaptor {
	readFile = func() ([]byte, error) {
		return []byte(`
Hardware        : BCM2708
Revision        : 0010
Serial          : 000000003bc748ea
`), nil
	}
	a := NewRaspiAdaptor("myAdaptor")
	a.Connect()
	return a
}

func TestRaspiAdaptor(t *testing.T) {
	readFile = func() ([]byte, error) {
		return []byte(`
Hardware        : BCM2708
Revision        : 0010
Serial          : 000000003bc748ea
`), nil
	}
	a := NewRaspiAdaptor("myAdaptor")
	gobot.Assert(t, a.Name(), "myAdaptor")
	gobot.Assert(t, a.i2cLocation, "/dev/i2c-1")
	gobot.Assert(t, a.revision, "3")

	readFile = func() ([]byte, error) {
		return []byte(`
Hardware        : BCM2708
Revision        : 000D
Serial          : 000000003bc748ea
`), nil
	}
	a = NewRaspiAdaptor("myAdaptor")
	gobot.Assert(t, a.i2cLocation, "/dev/i2c-1")
	gobot.Assert(t, a.revision, "2")

	readFile = func() ([]byte, error) {
		return []byte(`
Hardware        : BCM2708
Revision        : 0002
Serial          : 000000003bc748ea
`), nil
	}
	a = NewRaspiAdaptor("myAdaptor")
	gobot.Assert(t, a.i2cLocation, "/dev/i2c-0")
	gobot.Assert(t, a.revision, "1")

}
func TestRaspiAdaptorFinalize(t *testing.T) {
	a := initTestRaspiAdaptor()

	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/dev/pi-blaster",
		"/dev/i2c-1",
		"/dev/i2c-0",
	})

	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})

	a.DigitalWrite("3", 1)
	a.PwmWrite("7", 255)

	a.I2cStart(0xff)
	gobot.Assert(t, len(a.Finalize()), 0)
}

func TestRaspiAdaptorDigitalPWM(t *testing.T) {
	a := initTestRaspiAdaptor()

	gobot.Assert(t, a.PwmWrite("7", 4), nil)

	fs := sysfs.NewMockFilesystem([]string{
		"/dev/pi-blaster",
	})
	sysfs.SetFilesystem(fs)

	gobot.Assert(t, a.PwmWrite("7", 255), nil)

	gobot.Assert(t, strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0], "4=1")

	gobot.Assert(t, a.ServoWrite("11", 255), nil)

	gobot.Assert(t, strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0], "17=0.25")
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
	i, _ := a.DigitalRead("13")
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
	a.i2cDevice = &NullReadWriteCloser{}

	a.I2cWrite(0xff, []byte{0x00, 0x01})
	data, _ := a.I2cRead(0xff, 2)
	gobot.Assert(t, data, []byte{0x00, 0x01})
}
