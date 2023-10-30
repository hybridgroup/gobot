package tinkerboard

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
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

func TestName(t *testing.T) {
	a := NewAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Tinker Board"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestDigitalIO(t *testing.T) {
	// only basic tests needed, further tests are done in "digitalpinsadaptor.go"
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	_ = a.DigitalWrite("7", 1)
	assert.Equal(t, "1", fs.Files[gpio17Path+"value"].Contents)

	fs.Files[gpio160Path+"value"].Contents = "1"
	i, _ := a.DigitalRead("10")
	assert.Equal(t, 1, i)

	assert.ErrorContains(t, a.DigitalWrite("99", 1), "'99' is not a valid id for a digital pin")
	assert.NoError(t, a.Finalize())
}

func TestInvalidPWMPin(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	err := a.PwmWrite("666", 42)
	assert.ErrorContains(t, err, "'666' is not a valid id for a PWM pin")

	err = a.ServoWrite("666", 120)
	assert.ErrorContains(t, err, "'666' is not a valid id for a PWM pin")

	err = a.PwmWrite("3", 42)
	assert.ErrorContains(t, err, "'3' is not a valid id for a PWM pin")

	err = a.ServoWrite("3", 120)
	assert.ErrorContains(t, err, "'3' is not a valid id for a PWM pin")
}

func TestPwmWrite(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	err := a.PwmWrite("33", 100)
	assert.NoError(t, err)

	assert.Equal(t, "0", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwmEnablePath].Contents)
	assert.Equal(t, fmt.Sprintf("%d", 10000000), fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "3921568", fs.Files[pwmDutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwmPolarityPath].Contents)

	err = a.ServoWrite("33", 0)
	assert.NoError(t, err)

	assert.Equal(t, "500000", fs.Files[pwmDutyCyclePath].Contents)

	err = a.ServoWrite("33", 180)
	assert.NoError(t, err)

	assert.Equal(t, "2000000", fs.Files[pwmDutyCyclePath].Contents)
	assert.NoError(t, a.Finalize())
}

func TestSetPeriod(t *testing.T) {
	// arrange
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	newPeriod := uint32(2550000)
	// act
	err := a.SetPeriod("33", newPeriod)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, "0", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwmEnablePath].Contents)
	assert.Equal(t, fmt.Sprintf("%d", newPeriod), fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "0", fs.Files[pwmDutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwmPolarityPath].Contents)

	// arrange test for automatic adjustment of duty cycle to lower value
	err = a.PwmWrite("33", 127) // 127 is a little bit smaller than 50% of period
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%d", 1270000), fs.Files[pwmDutyCyclePath].Contents)
	newPeriod = newPeriod / 10

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%d", 127000), fs.Files[pwmDutyCyclePath].Contents)

	// arrange test for automatic adjustment of duty cycle to higher value
	newPeriod = newPeriod * 20

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%d", 2540000), fs.Files[pwmDutyCyclePath].Contents)
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	assert.NoError(t, a.DigitalWrite("7", 1))

	fs.WithWriteError = true

	err := a.Finalize()
	assert.Contains(t, err.Error(), "write error")
}

func TestFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	assert.NoError(t, a.PwmWrite("33", 1))

	fs.WithWriteError = true

	err := a.Finalize()
	assert.Contains(t, err.Error(), "write error")
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
	assert.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 4)
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
