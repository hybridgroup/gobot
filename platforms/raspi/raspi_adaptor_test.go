package raspi

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"runtime"
	"strconv"
	"sync"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/system"
)

// make sure that this Adaptor fulfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
var _ gobot.PWMPinnerProvider = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)
var _ spi.Connector = (*Adaptor)(nil)

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	a.Connect()
	return a, fs
}

func TestName(t *testing.T) {
	a := NewAdaptor()

	gobottest.Assert(t, strings.HasPrefix(a.Name(), "RaspberryPi"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestGetDefaultBus(t *testing.T) {
	const contentPattern = "Hardware        : BCM2708\n%sSerial          : 000000003bc748ea\n"
	var tests = map[string]struct {
		revisionPart string
		wantRev      string
		wantBus      int
	}{
		"no_revision": {
			wantRev: "0",
			wantBus: 0,
		},
		"rev_1": {
			revisionPart: "Revision        : 0002\n",
			wantRev:      "1",
			wantBus:      0,
		},
		"rev_2": {
			revisionPart: "Revision        : 000D\n",
			wantRev:      "2",
			wantBus:      1,
		},
		"rev_3": {
			revisionPart: "Revision        : 0010\n",
			wantRev:      "3",
			wantBus:      1,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			fs := a.sys.UseMockFilesystem([]string{infoFile})
			fs.Files[infoFile].Contents = fmt.Sprintf(contentPattern, tc.revisionPart)
			gobottest.Assert(t, a.revision, "")
			//act, will read and refresh the revision
			gotBus := a.GetDefaultBus()
			//assert
			gobottest.Assert(t, a.revision, tc.wantRev)
			gobottest.Assert(t, gotBus, tc.wantBus)
		})
	}
}

func TestFinalize(t *testing.T) {
	mockedPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/dev/pi-blaster",
		"/dev/i2c-1",
		"/dev/i2c-0",
		"/dev/spidev0.0",
		"/dev/spidev0.1",
	}
	a, _ := initTestAdaptorWithMockedFilesystem(mockedPaths)

	a.DigitalWrite("3", 1)
	a.PwmWrite("7", 255)

	a.GetConnection(0xff, 0)
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestDigitalPWM(t *testing.T) {
	mockedPaths := []string{"/dev/pi-blaster"}
	a, fs := initTestAdaptorWithMockedFilesystem(mockedPaths)
	a.PiBlasterPeriod = 20000000

	gobottest.Assert(t, a.PwmWrite("7", 4), nil)

	pin, _ := a.PWMPin("7")
	period, _ := pin.Period()
	gobottest.Assert(t, period, uint32(20000000))

	gobottest.Assert(t, a.PwmWrite("7", 255), nil)

	gobottest.Assert(t, strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0], "4=1")

	gobottest.Assert(t, a.ServoWrite("11", 90), nil)

	gobottest.Assert(t, strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0], "17=0.5")

	gobottest.Assert(t, a.PwmWrite("notexist", 1), errors.New("Not a valid pin"))
	gobottest.Assert(t, a.ServoWrite("notexist", 1), errors.New("Not a valid pin"))

	pin, _ = a.PWMPin("12")
	period, _ = pin.Period()
	gobottest.Assert(t, period, uint32(20000000))

	gobottest.Assert(t, pin.SetDutyCycle(1.5*1000*1000), nil)

	gobottest.Assert(t, strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0], "18=0.075")
}

func TestDigitalIO(t *testing.T) {
	mockedPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio4/value",
		"/sys/class/gpio/gpio4/direction",
		"/sys/class/gpio/gpio27/value",
		"/sys/class/gpio/gpio27/direction",
	}
	a, fs := initTestAdaptorWithMockedFilesystem(mockedPaths)

	err := a.DigitalWrite("7", 1)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio4/value"].Contents, "1")

	a.revision = "2"
	err = a.DigitalWrite("13", 1)
	gobottest.Assert(t, err, nil)

	i, err := a.DigitalRead("13")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("notexist", 1), errors.New("Not a valid pin"))

	fs.WithReadError = true
	_, err = a.DigitalRead("13")
	gobottest.Assert(t, err, errors.New("read error"))

	fs.WithWriteError = true
	_, err = a.DigitalRead("7")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestI2c(t *testing.T) {
	mockedPaths := []string{"/dev/i2c-1"}
	a, _ := initTestAdaptorWithMockedFilesystem(mockedPaths)
	a.sys.UseMockSyscall()

	con, err := a.GetConnection(0xff, 1)
	gobottest.Assert(t, err, nil)

	_, err = con.Write([]byte{0x00, 0x01})
	gobottest.Assert(t, err, nil)

	data := []byte{42, 42}
	_, err = con.Read(data)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	_, err = a.GetConnection(0xff, 51)
	gobottest.Assert(t, err, errors.New("Bus number 51 out of range"))

	a.revision = "0"
	gobottest.Assert(t, a.GetDefaultBus(), 0)

	a.revision = "2"
	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

func TestSPI(t *testing.T) {
	a := NewAdaptor()

	gobottest.Assert(t, a.GetSpiDefaultBus(), 0)
	gobottest.Assert(t, a.GetSpiDefaultChip(), 0)
	gobottest.Assert(t, a.GetSpiDefaultMode(), 0)
	gobottest.Assert(t, a.GetSpiDefaultMaxSpeed(), int64(500000))

	_, err := a.GetSpiConnection(10, 0, 0, 8, 500000)
	gobottest.Assert(t, err.Error(), "Bus number 10 out of range")

	// TODO: tests for real connection currently not possible, because not using system.Accessor using
	// TODO: test tx/rx here...
}

func TestDigitalPinConcurrency(t *testing.T) {

	oldProcs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(8)
	defer runtime.GOMAXPROCS(oldProcs)

	for retry := 0; retry < 20; retry++ {

		a := NewAdaptor()
		var wg sync.WaitGroup

		for i := 0; i < 20; i++ {
			wg.Add(1)
			pinAsString := strconv.Itoa(i)
			go func(pin string) {
				defer wg.Done()
				a.DigitalPin(pin)
			}(pinAsString)
		}

		wg.Wait()
	}
}

func TestPWMPin(t *testing.T) {
	a := NewAdaptor()

	gobottest.Assert(t, len(a.pwmPins), 0)

	a.revision = "3"
	firstSysPin, err := a.PWMPin("35")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pwmPins), 1)

	secondSysPin, err := a.PWMPin("35")

	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pwmPins), 1)
	gobottest.Assert(t, firstSysPin, secondSysPin)

	otherSysPin, err := a.PWMPin("36")

	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pwmPins), 2)
	gobottest.Refute(t, firstSysPin, otherSysPin)
}
