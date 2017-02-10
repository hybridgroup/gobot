package raspi

import (
	"strings"
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

func initTestAdaptor() *Adaptor {
	readFile = func() ([]byte, error) {
		return []byte(`
Hardware        : BCM2708
Revision        : 0010
Serial          : 000000003bc748ea
`), nil
	}
	a := NewAdaptor()
	a.Connect()
	return a
}

func TestAdaptor(t *testing.T) {
	readFile = func() ([]byte, error) {
		return []byte(`
Hardware        : BCM2708
Revision        : 0010
Serial          : 000000003bc748ea
`), nil
	}
	a := NewAdaptor()
	gobottest.Assert(t, a.Name(), "RaspberryPi")
	gobottest.Assert(t, a.i2cDefaultBus, 1)
	gobottest.Assert(t, a.revision, "3")

	readFile = func() ([]byte, error) {
		return []byte(`
Hardware        : BCM2708
Revision        : 000D
Serial          : 000000003bc748ea
`), nil
	}
	a = NewAdaptor()
	gobottest.Assert(t, a.i2cDefaultBus, 1)
	gobottest.Assert(t, a.revision, "2")

	readFile = func() ([]byte, error) {
		return []byte(`
Hardware        : BCM2708
Revision        : 0002
Serial          : 000000003bc748ea
`), nil
	}
	a = NewAdaptor()
	gobottest.Assert(t, a.i2cDefaultBus, 0)
	gobottest.Assert(t, a.revision, "1")

}
func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()

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

	a.GetConnection(0xff, 0)
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestAdaptorDigitalPWM(t *testing.T) {
	a := initTestAdaptor()

	gobottest.Assert(t, a.PwmWrite("7", 4), nil)

	fs := sysfs.NewMockFilesystem([]string{
		"/dev/pi-blaster",
	})
	sysfs.SetFilesystem(fs)

	gobottest.Assert(t, a.PwmWrite("7", 255), nil)

	gobottest.Assert(t, strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0], "4=1")

	gobottest.Assert(t, a.ServoWrite("11", 255), nil)

	gobottest.Assert(t, strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0], "17=0.25")
}

func TestAdaptorDigitalIO(t *testing.T) {
	a := initTestAdaptor()
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
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio4/value"].Contents, "1")

	a.DigitalWrite("13", 1)
	i, _ := a.DigitalRead("13")
	gobottest.Assert(t, i, 1)
}

func TestAdaptorI2c(t *testing.T) {
	a := initTestAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})

	con, err := a.GetConnection(0xff, 1)
	gobottest.Assert(t, err, nil)

	con.Write([]byte{0x00, 0x01})
	data := []byte{42, 42}
	con.Read(data)
	gobottest.Assert(t, data, []byte{0x00, 0x01})
}
