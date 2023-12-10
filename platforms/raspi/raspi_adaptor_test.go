package raspi

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/drivers/spi"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	pwmDir           = "/sys/class/pwm/pwmchip0/" //nolint:gosec // false positive
	pwmPwmDir        = pwmDir + "pwm0/"
	pwmExportPath    = pwmDir + "export"
	pwmUnexportPath  = pwmDir + "unexport"
	pwmEnablePath    = pwmPwmDir + "enable"
	pwmPeriodPath    = pwmPwmDir + "period"
	pwmDutyCyclePath = pwmPwmDir + "duty_cycle"
	pwmPolarityPath  = pwmPwmDir + "polarity"

	pwmInvertedIdentifier = "inversed"
)

var pwmMockPaths = []string{
	pwmExportPath,
	pwmUnexportPath,
	pwmEnablePath,
	pwmPeriodPath,
	pwmDutyCyclePath,
	pwmPolarityPath,
}

// make sure that this Adaptor fulfills all the required interfaces
var (
	_ gobot.Adaptor               = (*Adaptor)(nil)
	_ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
	_ gobot.PWMPinnerProvider     = (*Adaptor)(nil)
	_ gpio.DigitalReader          = (*Adaptor)(nil)
	_ gpio.DigitalWriter          = (*Adaptor)(nil)
	_ gpio.PwmWriter              = (*Adaptor)(nil)
	_ gpio.ServoWriter            = (*Adaptor)(nil)
	_ aio.AnalogReader            = (*Adaptor)(nil)
	_ i2c.Connector               = (*Adaptor)(nil)
	_ spi.Connector               = (*Adaptor)(nil)
)

func preparePwmFs(fs *system.MockFilesystem) {
	fs.Files[pwmEnablePath].Contents = "0"
	fs.Files[pwmPeriodPath].Contents = "0"
	fs.Files[pwmDutyCyclePath].Contents = "0"
	fs.Files[pwmPolarityPath].Contents = pwmInvertedIdentifier
}

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
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
	require.NoError(t, a.Finalize())
}

func TestAnalog(t *testing.T) {
	mockPaths := []string{
		"/sys/class/thermal/thermal_zone0/temp",
	}

	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)

	fs.Files["/sys/class/thermal/thermal_zone0/temp"].Contents = "567\n"
	got, err := a.AnalogRead("thermal_zone0")
	require.NoError(t, err)
	assert.Equal(t, 567, got)

	_, err = a.AnalogRead("thermal_zone10")
	require.ErrorContains(t, err, "'thermal_zone10' is not a valid id for a analog pin")

	fs.WithReadError = true
	_, err = a.AnalogRead("thermal_zone0")
	require.ErrorContains(t, err, "read error")
	fs.WithReadError = false

	require.NoError(t, a.Finalize())
}

func TestPwmWrite(t *testing.T) {
	// arrange
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)
	// act
	err := a.PwmWrite("pwm0", 100)
	// assert
	require.NoError(t, err)
	assert.Equal(t, "0", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwmEnablePath].Contents)
	assert.Equal(t, "10000000", fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "3921568", fs.Files[pwmDutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwmPolarityPath].Contents)
	// act & assert invalid pin
	err = a.PwmWrite("pwm1", 42)
	require.ErrorContains(t, err, "'pwm1' is not a valid pin id for raspi revision 0")
	require.NoError(t, a.Finalize())
}

func TestServoWrite(t *testing.T) {
	// arrange: prepare 50Hz for servos
	const (
		pin         = "pwm0"
		fiftyHzNano = 20000000
	)
	a := NewAdaptor(adaptors.WithPWMDefaultPeriodForPin(pin, fiftyHzNano))
	fs := a.sys.UseMockFilesystem(pwmMockPaths)
	preparePwmFs(fs)
	require.NoError(t, a.Connect())
	// act & assert for 0° (min default value)
	err := a.ServoWrite(pin, 0)
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(fiftyHzNano), fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "500000", fs.Files[pwmDutyCyclePath].Contents)
	// act & assert for 180° (max default value)
	err = a.ServoWrite(pin, 180)
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(fiftyHzNano), fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "2500000", fs.Files[pwmDutyCyclePath].Contents)
	// act & assert invalid pins
	err = a.ServoWrite("3", 120)
	require.ErrorContains(t, err, "'3' is not a valid pin id for raspi revision 0")
	require.NoError(t, a.Finalize())
}

func TestPWMWrite_piPlaster(t *testing.T) {
	// arrange
	const hundredHzNano = 10000000
	mockedPaths := []string{"/dev/pi-blaster"}
	a := NewAdaptor(adaptors.WithPWMUsePiBlaster())
	fs := a.sys.UseMockFilesystem(mockedPaths)
	require.NoError(t, a.Connect())
	// act & assert: Write & Pin & Period
	require.NoError(t, a.PwmWrite("7", 255))
	assert.Equal(t, "4=1", strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0])
	pin, _ := a.PWMPin("7")
	period, err := pin.Period()
	require.NoError(t, err)
	assert.Equal(t, uint32(hundredHzNano), period)
	// act & assert: nonexistent pin
	require.ErrorContains(t, a.PwmWrite("notexist", 1), "'notexist' is not a valid pin id for raspi revision 0")
	// act & assert: SetDutyCycle
	pin, _ = a.PWMPin("12")
	require.NoError(t, pin.SetDutyCycle(1.5*1000*1000))
	assert.Equal(t, "18=0.15", strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0])
}

func TestPWM_piPlaster(t *testing.T) {
	// arrange
	const fiftyHzNano = 20000000 // 20 ms
	mockedPaths := []string{"/dev/pi-blaster"}
	a := NewAdaptor(adaptors.WithPWMUsePiBlaster(), adaptors.WithPWMDefaultPeriod(fiftyHzNano))
	fs := a.sys.UseMockFilesystem(mockedPaths)
	require.NoError(t, a.Connect())
	// act & assert: Pin & Period
	pin, _ := a.PWMPin("7")
	period, err := pin.Period()
	require.NoError(t, err)
	assert.Equal(t, uint32(fiftyHzNano), period)
	// act & assert for 180° (max default value), 2.5 ms => 12.5%
	require.NoError(t, a.ServoWrite("11", 180))
	assert.Equal(t, "17=0.125", strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0])
	// act & assert for 90° (center value), 1.5 ms => 7.5% duty
	require.NoError(t, a.ServoWrite("11", 90))
	assert.Equal(t, "17=0.075", strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0])
	// act & assert for 0° (min default value), 0.5 ms => 2.5% duty
	require.NoError(t, a.ServoWrite("11", 0))
	assert.Equal(t, "17=0.025", strings.Split(fs.Files["/dev/pi-blaster"].Contents, "\n")[0])
	// act & assert: nonexistent pin
	require.ErrorContains(t, a.ServoWrite("notexist", 1), "'notexist' is not a valid pin id for raspi revision 0")
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
	require.NoError(t, err)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio4/value"].Contents)

	a.revision = "2"
	err = a.DigitalWrite("13", 1)
	require.NoError(t, err)

	i, err := a.DigitalRead("13")
	require.NoError(t, err)
	assert.Equal(t, 1, i)

	require.ErrorContains(t, a.DigitalWrite("notexist", 1), "'notexist' is not a valid pin id for raspi revision 2")
	require.NoError(t, a.Finalize())
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
	require.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 1)
	require.NoError(t, err)
	_, err = con.Write([]byte{0xbf})
	require.NoError(t, err)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	require.ErrorContains(t, err, "close error")
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

func Test_translateAnalogPin(t *testing.T) {
	mockedPaths := []string{
		"/sys/class/thermal/thermal_zone0/temp",
		"/sys/class/thermal/thermal_zone1/temp",
	}
	tests := map[string]struct {
		id           string
		wantPath     string
		wantReadable bool
		wantBufLen   uint16
		wantErr      string
	}{
		"translate_thermal_zone0": {
			id:           "thermal_zone0",
			wantPath:     "/sys/class/thermal/thermal_zone0/temp",
			wantReadable: true,
			wantBufLen:   7,
		},
		"unknown_id": {
			id:      "thermal_zone1",
			wantErr: "'thermal_zone1' is not a valid id for a analog pin",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a, _ := initTestAdaptorWithMockedFilesystem(mockedPaths)
			// act
			path, r, w, buf, err := a.translateAnalogPin(tc.id)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantPath, path)
			assert.Equal(t, tc.wantReadable, r)
			assert.False(t, w)
			assert.Equal(t, tc.wantBufLen, buf)
		})
	}
}

func Test_getPinTranslatorFunction(t *testing.T) {
	tests := map[string]struct {
		id       string
		revision string
		wantPath string
		wantLine int
		wantErr  string
	}{
		"translate_12_rev0": {
			id:       "12",
			wantPath: "gpiochip0",
			wantLine: 18,
		},
		"translate_13_rev0": {
			id:      "13",
			wantErr: "'13' is not a valid pin id for raspi revision 0",
		},
		"translate_13_rev1": {
			id:       "13",
			revision: "1",
			wantPath: "gpiochip0",
			wantLine: 21,
		},
		"translate_29_rev1": {
			id:       "29",
			revision: "1",
			wantErr:  "'29' is not a valid pin id for raspi revision 1",
		},
		"translate_29_rev3": {
			id:       "29",
			revision: "3",
			wantPath: "gpiochip0",
			wantLine: 5,
		},
		"translate_pwm0_rev0": {
			id:       "pwm0",
			wantPath: "/sys/class/pwm/pwmchip0",
			wantLine: 0,
		},
		"translate_pwm1_rev0": {
			id:      "pwm1",
			wantErr: "'pwm1' is not a valid pin id for raspi revision 0",
		},
		"translate_pwm1_rev3": {
			id:       "pwm1",
			revision: "3",
			wantPath: "/sys/class/pwm/pwmchip0",
			wantLine: 1,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			a.revision = tc.revision
			// act
			f := a.getPinTranslatorFunction()
			path, line, err := f(tc.id)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantPath, path)
			assert.Equal(t, tc.wantLine, line)
		})
	}
}
