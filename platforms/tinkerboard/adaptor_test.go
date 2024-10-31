package tinkerboard

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	gpio17Path  = "/sys/class/gpio/gpio17/"
	gpio160Path = "/sys/class/gpio/gpio160/"
)

const (
	pwmDir           = "/sys/devices/platform/ff680020.pwm/pwm/pwmchip2/" //nolint:gosec // false positive
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

var gpioMockPaths = []string{
	"/sys/class/gpio/export",
	"/sys/class/gpio/unexport",
	gpio17Path + "value",
	gpio17Path + "direction",
	gpio160Path + "value",
	gpio160Path + "direction",
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

func TestNewAdaptor(t *testing.T) {
	// arrange & act
	a := NewAdaptor()
	// assert
	assert.IsType(t, &Adaptor{}, a)
	assert.True(t, strings.HasPrefix(a.Name(), "Tinker Board"))
	assert.NotNil(t, a.sys)
	assert.NotNil(t, a.mutex)
	assert.NotNil(t, a.AnalogPinsAdaptor)
	assert.NotNil(t, a.DigitalPinsAdaptor)
	assert.NotNil(t, a.PWMPinsAdaptor)
	assert.NotNil(t, a.I2cBusAdaptor)
	assert.NotNil(t, a.SpiBusAdaptor)
	assert.NotNil(t, a.OneWireBusAdaptor)
	// act & assert
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestNewAdaptorWithOption(t *testing.T) {
	// arrange & act
	a := NewAdaptor(adaptors.WithGpiosActiveLow("1"))
	// assert
	require.NoError(t, a.Connect())
}

func TestDigitalIO(t *testing.T) {
	// only basic tests needed, further tests are done in "digitalpinsadaptor.go"
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	_ = a.DigitalWrite("7", 1)
	assert.Equal(t, "1", fs.Files[gpio17Path+"value"].Contents)

	fs.Files[gpio160Path+"value"].Contents = "1"
	i, _ := a.DigitalRead("10")
	assert.Equal(t, 1, i)

	require.ErrorContains(t, a.DigitalWrite("99", 1), "'99' is not a valid id for a digital pin")
	require.NoError(t, a.Finalize())
}

func TestAnalogRead(t *testing.T) {
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
	err := a.PwmWrite("33", 100)
	// assert
	require.NoError(t, err)
	assert.Equal(t, "0", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwmEnablePath].Contents)
	assert.Equal(t, "10000000", fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "3921568", fs.Files[pwmDutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwmPolarityPath].Contents)
	// act & assert invalid pin
	err = a.PwmWrite("666", 42)
	require.ErrorContains(t, err, "'666' is not a valid id for a PWM pin")

	require.NoError(t, a.Finalize())
}

func TestServoWrite(t *testing.T) {
	// arrange: prepare 50Hz for servos
	const (
		pin         = "33"
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
	require.ErrorContains(t, err, "'3' is not a valid id for a PWM pin")

	require.NoError(t, a.Finalize())
}

func TestSetPeriod(t *testing.T) {
	// arrange
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	newPeriod := uint32(2550000)
	// act
	err := a.SetPeriod("33", newPeriod)
	// assert
	require.NoError(t, err)
	assert.Equal(t, "0", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwmEnablePath].Contents)
	assert.Equal(t, fmt.Sprintf("%d", newPeriod), fs.Files[pwmPeriodPath].Contents) //nolint:perfsprint // ok here
	assert.Equal(t, "0", fs.Files[pwmDutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwmPolarityPath].Contents)

	// arrange test for automatic adjustment of duty cycle to lower value
	err = a.PwmWrite("33", 127) // 127 is a little bit smaller than 50% of period
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(1270000), fs.Files[pwmDutyCyclePath].Contents)
	newPeriod = newPeriod / 10

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(127000), fs.Files[pwmDutyCyclePath].Contents)

	// arrange test for automatic adjustment of duty cycle to higher value
	newPeriod = newPeriod * 20

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(2540000), fs.Files[pwmDutyCyclePath].Contents)
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	require.NoError(t, a.DigitalWrite("7", 1))

	fs.WithWriteError = true

	err := a.Finalize()
	require.ErrorContains(t, err, "write error")
}

func TestFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	require.NoError(t, a.PwmWrite("33", 1))

	fs.WithWriteError = true

	err := a.Finalize()
	require.ErrorContains(t, err, "write error")
}

func TestSpiDefaultValues(t *testing.T) {
	a := NewAdaptor()

	assert.Equal(t, 0, a.SpiDefaultBusNumber())
	assert.Equal(t, 0, a.SpiDefaultChipNumber())
	assert.Equal(t, 0, a.SpiDefaultMode())
	assert.Equal(t, 8, a.SpiDefaultBitCount())
	assert.Equal(t, int64(500000), a.SpiDefaultMaxSpeed())
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	assert.Equal(t, 1, a.DefaultI2cBus())
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-4"})
	require.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 4)
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
		"number_1_error": {
			busNr:   1,
			wantErr: fmt.Errorf("Bus number 1 out of range"),
		},
		"number_2_ok": {
			busNr: 2,
		},
		"number_3_error": {
			busNr:   3,
			wantErr: fmt.Errorf("Bus number 3 out of range"),
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
		"number_2_ok": {
			busNr: 2,
		},
		"number_3_ok": {
			busNr: 3,
		},
		"number_4_ok": {
			busNr: 4,
		},
		"number_5_error": {
			busNr:   5,
			wantErr: fmt.Errorf("Bus number 5 out of range"),
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

func Test_translateDigitalPin(t *testing.T) {
	tests := map[string]struct {
		access   string
		pin      string
		wantChip string
		wantLine int
		wantErr  error
	}{
		"cdev_ok": {
			access:   "cdev",
			pin:      "7",
			wantChip: "gpiochip0",
			wantLine: 17,
		},
		"sysfs_ok": {
			access:   "sysfs",
			pin:      "7",
			wantChip: "",
			wantLine: 17,
		},
		"unknown_pin": {
			pin:      "99",
			wantChip: "",
			wantLine: -1,
			wantErr:  fmt.Errorf("'99' is not a valid id for a digital pin"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			a.sys.UseDigitalPinAccessWithMockFs(tc.access, []string{})
			// act
			chip, line, err := a.translateDigitalPin(tc.pin)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantChip, chip)
			assert.Equal(t, tc.wantLine, line)
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
		"translate_thermal_zone1": {
			id:           "thermal_zone1",
			wantPath:     "/sys/class/thermal/thermal_zone1/temp",
			wantReadable: true,
			wantBufLen:   7,
		},
		"unknown_id": {
			id:      "99",
			wantErr: "'99' is not a valid id for a analog pin",
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

func Test_translatePWMPin(t *testing.T) {
	basePaths := []string{
		"/sys/devices/platform/ff680020.pwm/pwm/",
		"/sys/devices/platform/ff680030.pwm/pwm/",
	}
	tests := map[string]struct {
		pin         string
		chip        string
		wantDir     string
		wantChannel int
		wantErr     error
	}{
		"32_chip0": {
			pin:         "32",
			chip:        "pwmchip0",
			wantDir:     "/sys/devices/platform/ff680030.pwm/pwm/pwmchip0",
			wantChannel: 0,
		},
		"32_chip1": {
			pin:         "32",
			chip:        "pwmchip1",
			wantDir:     "/sys/devices/platform/ff680030.pwm/pwm/pwmchip1",
			wantChannel: 0,
		},
		"32_chip2": {
			pin:         "32",
			chip:        "pwmchip2",
			wantDir:     "/sys/devices/platform/ff680030.pwm/pwm/pwmchip2",
			wantChannel: 0,
		},
		"32_chip3": {
			pin:         "32",
			chip:        "pwmchip3",
			wantDir:     "/sys/devices/platform/ff680030.pwm/pwm/pwmchip3",
			wantChannel: 0,
		},
		"33_chip0": {
			pin:         "33",
			chip:        "pwmchip0",
			wantDir:     "/sys/devices/platform/ff680020.pwm/pwm/pwmchip0",
			wantChannel: 0,
		},
		"33_chip1": {
			pin:         "33",
			chip:        "pwmchip1",
			wantDir:     "/sys/devices/platform/ff680020.pwm/pwm/pwmchip1",
			wantChannel: 0,
		},
		"33_chip2": {
			pin:         "33",
			chip:        "pwmchip2",
			wantDir:     "/sys/devices/platform/ff680020.pwm/pwm/pwmchip2",
			wantChannel: 0,
		},
		"invalid_pin": {
			pin:         "7",
			wantDir:     "",
			wantChannel: -1,
			wantErr:     fmt.Errorf("'7' is not a valid id for a PWM pin"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			mockedPaths := []string{}
			for _, base := range basePaths {
				mockedPaths = append(mockedPaths, base+tc.chip+"/")
			}
			a, _ := initTestAdaptorWithMockedFilesystem(mockedPaths)
			// act
			dir, channel, err := a.translatePWMPin(tc.pin)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantDir, dir)
			assert.Equal(t, tc.wantChannel, channel)
		})
	}
}
