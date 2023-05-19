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
			gotBus := a.DefaultI2cBus()
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

	a.GetI2cConnection(0xff, 0)
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
	gobottest.Assert(t, a.Finalize(), nil)
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
	if err := a.Connect(); err != nil {
		panic(err)
	}

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

func TestPWMPinsReConnect(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.revision = "3"
	if err := a.Connect(); err != nil {
		panic(err)
	}

	_, err := a.PWMPin("35")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pwmPins), 1)
	gobottest.Assert(t, a.Finalize(), nil)
	// act
	err = a.Connect()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pwmPins), 0)
	_, _ = a.PWMPin("35")
	_, err = a.PWMPin("36")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pwmPins), 2)
}

func TestSpiDefaultValues(t *testing.T) {
	a := NewAdaptor()

	gobottest.Assert(t, a.SpiDefaultBusNumber(), 0)
	gobottest.Assert(t, a.SpiDefaultChipNumber(), 0)
	gobottest.Assert(t, a.SpiDefaultMode(), 0)
	gobottest.Assert(t, a.SpiDefaultMaxSpeed(), int64(500000))
}

func TestI2cDefaultBus(t *testing.T) {
	mockedPaths := []string{"/dev/i2c-1"}
	a, _ := initTestAdaptorWithMockedFilesystem(mockedPaths)
	a.sys.UseMockSyscall()

	a.revision = "0"
	gobottest.Assert(t, a.DefaultI2cBus(), 0)

	a.revision = "2"
	gobottest.Assert(t, a.DefaultI2cBus(), 1)
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-1"})
	gobottest.Assert(t, a.Connect(), nil)
	con, err := a.GetI2cConnection(0xff, 1)
	gobottest.Assert(t, err, nil)
	_, err = con.Write([]byte{0xbf})
	gobottest.Assert(t, err, nil)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	gobottest.Assert(t, strings.Contains(err.Error(), "close error"), true)
}

func Test_validateSpiBusNumber(t *testing.T) {
	var tests = map[string]struct {
		busNr   int
		wantErr error
	}{
		"number_negative_error": {
			busNr:   -1,
			wantErr: fmt.Errorf("Bus number -1 out of range"),
		},
		"number_0_ok": {
			busNr: 0,
		},
		"number_1_ok": {
			busNr: 1,
		},
		"number_2_error": {
			busNr:   2,
			wantErr: fmt.Errorf("Bus number 2 out of range"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			// act
			err := a.validateSpiBusNumber(tc.busNr)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
		})
	}
}

func Test_validateI2cBusNumber(t *testing.T) {
	var tests = map[string]struct {
		busNr   int
		wantErr error
	}{
		"number_negative_error": {
			busNr:   -1,
			wantErr: fmt.Errorf("Bus number -1 out of range"),
		},
		"number_0_ok": {
			busNr: 0,
		},
		"number_1_ok": {
			busNr: 1,
		},
		"number_2_not_ok": {
			busNr:   2,
			wantErr: fmt.Errorf("Bus number 2 out of range"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			// act
			err := a.validateI2cBusNumber(tc.busNr)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
		})
	}
}
