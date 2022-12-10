package tinkerboard

import (
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/system"
)

const (
	gpio17Path  = "/sys/class/gpio/gpio17/"
	gpio160Path = "/sys/class/gpio/gpio160/"
)

const (
	pwmDir           = "/sys/devices/platform/ff680020.pwm/pwm/pwmchip2/"
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
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
var _ gobot.PWMPinnerProvider = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

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
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Tinker Board"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestDigitalIO(t *testing.T) {
	// only basic tests needed, further tests are done in "digitalpinsadaptor.go"
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	a.DigitalWrite("7", 1)
	gobottest.Assert(t, fs.Files[gpio17Path+"value"].Contents, "1")

	fs.Files[gpio160Path+"value"].Contents = "1"
	i, _ := a.DigitalRead("10")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("99", 1), fmt.Errorf("'99' is not a valid id for a digital pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestInvalidPWMPin(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	err := a.PwmWrite("666", 42)
	gobottest.Assert(t, err.Error(), "'666' is not a valid id for a PWM pin")

	err = a.ServoWrite("666", 120)
	gobottest.Assert(t, err.Error(), "'666' is not a valid id for a PWM pin")

	err = a.PwmWrite("3", 42)
	gobottest.Assert(t, err.Error(), "'3' is not a valid id for a PWM pin")

	err = a.ServoWrite("3", 120)
	gobottest.Assert(t, err.Error(), "'3' is not a valid id for a PWM pin")
}

func TestPwmWrite(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	err := a.PwmWrite("33", 100)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files[pwmExportPath].Contents, "0")
	gobottest.Assert(t, fs.Files[pwmEnablePath].Contents, "1")
	gobottest.Assert(t, fs.Files[pwmPeriodPath].Contents, fmt.Sprintf("%d", 10000000))
	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, "3921568")
	gobottest.Assert(t, fs.Files[pwmPolarityPath].Contents, "normal")

	err = a.ServoWrite("33", 0)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, "500000")

	err = a.ServoWrite("33", 180)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, "2000000")
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestSetPeriod(t *testing.T) {
	// arrange
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	newPeriod := uint32(2550000)
	// act
	err := a.SetPeriod("33", newPeriod)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwmExportPath].Contents, "0")
	gobottest.Assert(t, fs.Files[pwmEnablePath].Contents, "1")
	gobottest.Assert(t, fs.Files[pwmPeriodPath].Contents, fmt.Sprintf("%d", newPeriod))
	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, "0")
	gobottest.Assert(t, fs.Files[pwmPolarityPath].Contents, "normal")

	// arrange test for automatic adjustment of duty cycle to lower value
	err = a.PwmWrite("33", 127) // 127 is a little bit smaller than 50% of period
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, fmt.Sprintf("%d", 1270000))
	newPeriod = newPeriod / 10

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, fmt.Sprintf("%d", 127000))

	// arrange test for automatic adjustment of duty cycle to higher value
	newPeriod = newPeriod * 20

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, fmt.Sprintf("%d", 2540000))
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	gobottest.Assert(t, a.DigitalWrite("7", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	gobottest.Assert(t, a.PwmWrite("33", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestSpiDefaultValues(t *testing.T) {
	a := NewAdaptor()

	gobottest.Assert(t, a.SpiDefaultBusNumber(), 0)
	gobottest.Assert(t, a.SpiDefaultChipNumber(), 0)
	gobottest.Assert(t, a.SpiDefaultMode(), 0)
	gobottest.Assert(t, a.SpiDefaultBitCount(), 8)
	gobottest.Assert(t, a.SpiDefaultMaxSpeed(), int64(500000))
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, a.DefaultI2cBus(), 1)
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-4"})
	gobottest.Assert(t, a.Connect(), nil)
	con, err := a.GetI2cConnection(0xff, 4)
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
			gobottest.Assert(t, err, tc.wantErr)
		})
	}
}

func Test_translateDigitalPin(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, err, tc.wantErr)
			gobottest.Assert(t, chip, tc.wantChip)
			gobottest.Assert(t, line, tc.wantLine)
		})
	}
}

func Test_translatePWMPin(t *testing.T) {
	basePaths := []string{
		"/sys/devices/platform/ff680020.pwm/pwm/",
		"/sys/devices/platform/ff680030.pwm/pwm/",
	}
	var tests = map[string]struct {
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
			gobottest.Assert(t, err, tc.wantErr)
			gobottest.Assert(t, dir, tc.wantDir)
			gobottest.Assert(t, channel, tc.wantChannel)
		})
	}
}
