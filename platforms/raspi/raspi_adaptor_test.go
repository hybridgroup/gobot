package raspi

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/drivers/spi"
	"gobot.io/x/gobot/v2/system"
)

// make sure that this Adaptor fulfills all the required interfaces
var (
	_ gobot.Adaptor               = (*Adaptor)(nil)
	_ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
	_ gobot.PWMPinnerProvider     = (*Adaptor)(nil)
	_ gpio.DigitalReader          = (*Adaptor)(nil)
	_ gpio.DigitalWriter          = (*Adaptor)(nil)
	_ gpio.PwmWriter              = (*Adaptor)(nil)
	_ gpio.ServoWriter            = (*Adaptor)(nil)
	_ i2c.Connector               = (*Adaptor)(nil)
	_ spi.Connector               = (*Adaptor)(nil)
)

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	_ = a.Connect()
	return a, fs
}

func TestName(t *testing.T) {
	a := NewAdaptor()

	assert.True(t, strings.HasPrefix(a.Name(), "RaspberryPi"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestGetDefaultBus(t *testing.T) {
	const contentPattern = "Hardware        : BCM2708\n%sSerial          : 000000003bc748ea\n"
	tests := map[string]struct {
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
			assert.Equal(t, "", a.revision)
			// act, will read and refresh the revision
			gotBus := a.DefaultI2cBus()
			// assert
			assert.Equal(t, tc.wantRev, a.revision)
			assert.Equal(t, tc.wantBus, gotBus)
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

	_ = a.DigitalWrite("3", 1)
	_ = a.PwmWrite("7", 255)

	_, _ = a.GetI2cConnection(0xff, 0)
	assert.NoError(t, a.Finalize())
}

func TestDigitalPWM(t *testing.T) {
	mockedPaths := []string{"/dev/pi-blaster"}
	a, fs := initTestAdaptorWithMockedFilesystem(mockedPaths)
	a.PiBlasterPeriod = 20000000

	assert.NoError(t, a.PwmWrite("7", 4))

	pin, _ := a.PWMPin("7")
	period, _ := pin.Period()
	assert.Equal(t, uint32(20000000), period)

	assert.NoError(t, a.PwmWrite("7", 255))

	assert.Equal(t, "4=1", strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0])

	assert.NoError(t, a.ServoWrite("11", 90))

	assert.Equal(t, "17=0.5", strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0])

	assert.ErrorContains(t, a.PwmWrite("notexist", 1), "Not a valid pin")
	assert.ErrorContains(t, a.ServoWrite("notexist", 1), "Not a valid pin")

	pin, _ = a.PWMPin("12")
	period, _ = pin.Period()
	assert.Equal(t, uint32(20000000), period)

	assert.NoError(t, pin.SetDutyCycle(1.5*1000*1000))

	assert.Equal(t, "18=0.075", strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0])
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
	assert.NoError(t, err)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio4/value"].Contents)

	a.revision = "2"
	err = a.DigitalWrite("13", 1)
	assert.NoError(t, err)

	i, err := a.DigitalRead("13")
	assert.NoError(t, err)
	assert.Equal(t, 1, i)

	assert.ErrorContains(t, a.DigitalWrite("notexist", 1), "Not a valid pin")
	assert.NoError(t, a.Finalize())
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
				_, _ = a.DigitalPin(pin)
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

	assert.Equal(t, 0, len(a.pwmPins))

	a.revision = "3"
	firstSysPin, err := a.PWMPin("35")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(a.pwmPins))

	secondSysPin, err := a.PWMPin("35")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(a.pwmPins))
	assert.Equal(t, secondSysPin, firstSysPin)

	otherSysPin, err := a.PWMPin("36")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(a.pwmPins))
	assert.NotEqual(t, otherSysPin, firstSysPin)
}

func TestPWMPinsReConnect(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.revision = "3"
	if err := a.Connect(); err != nil {
		panic(err)
	}

	_, err := a.PWMPin("35")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(a.pwmPins))
	assert.NoError(t, a.Finalize())
	// act
	err = a.Connect()
	// assert
	assert.NoError(t, err)
	assert.Equal(t, 0, len(a.pwmPins))
	_, _ = a.PWMPin("35")
	_, err = a.PWMPin("36")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(a.pwmPins))
}

func TestSpiDefaultValues(t *testing.T) {
	a := NewAdaptor()

	assert.Equal(t, 0, a.SpiDefaultBusNumber())
	assert.Equal(t, 0, a.SpiDefaultChipNumber())
	assert.Equal(t, 0, a.SpiDefaultMode())
	assert.Equal(t, int64(500000), a.SpiDefaultMaxSpeed())
}

func TestI2cDefaultBus(t *testing.T) {
	mockedPaths := []string{"/dev/i2c-1"}
	a, _ := initTestAdaptorWithMockedFilesystem(mockedPaths)
	a.sys.UseMockSyscall()

	a.revision = "0"
	assert.Equal(t, 0, a.DefaultI2cBus())

	a.revision = "2"
	assert.Equal(t, 1, a.DefaultI2cBus())
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-1"})
	assert.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 1)
	assert.NoError(t, err)
	_, err = con.Write([]byte{0xbf})
	assert.NoError(t, err)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	assert.Contains(t, err.Error(), "close error")
}

func Test_validateSpiBusNumber(t *testing.T) {
	tests := map[string]struct {
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
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func Test_validateI2cBusNumber(t *testing.T) {
	tests := map[string]struct {
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
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
