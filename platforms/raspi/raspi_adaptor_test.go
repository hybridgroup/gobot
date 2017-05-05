package raspi

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
	"strconv"
	"sync"
	"runtime"
)

// make sure that this Adaptor fullfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

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

func TestRaspiAdaptorName(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "RaspberryPi"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
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
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "RaspberryPi"), true)
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

	gobottest.Assert(t, a.PwmWrite("notexist", 1), errors.New("Not a valid pin"))
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

	gobottest.Assert(t, a.DigitalWrite("notexist", 1), errors.New("Not a valid pin"))
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

	_, err = a.GetConnection(0xff, 51)
	gobottest.Assert(t, err, errors.New("Bus number 51 out of range"))

	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

// package internals testing

func TestAdaptor_concurrentDigitalPin(t *testing.T) {

	oldProcs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(8)

	for retry := 0; retry < 20; retry++ {

		a := initTestAdaptor()
		var wg sync.WaitGroup
		wg.Add(20)

		for i := 0; i < 20; i++ {
			pinAsString := strconv.Itoa(i)
			go func() {
				defer wg.Done()
				a.digitalPin(pinAsString, sysfs.IN)
			}()
		}

		wg.Wait()
	}

	runtime.GOMAXPROCS(oldProcs)

}

func TestAdaptor_pwmPin(t *testing.T) {
	a := initTestAdaptor()

	gobottest.Assert(t, len(a.pwmPins), 0)

	firstSysPin, err := a.pwmPin("35")

	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pwmPins), 1)

	secondSysPin, err := a.pwmPin("35")

	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pwmPins), 1)
	gobottest.Assert(t, firstSysPin, secondSysPin)

	otherSysPin, err := a.pwmPin("36")

	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pwmPins), 2)
	gobottest.Refute(t, firstSysPin, otherSysPin)
}
