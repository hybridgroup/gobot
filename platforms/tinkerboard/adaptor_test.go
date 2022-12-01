package tinkerboard

import (
	"errors"
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
	pwm2Dir           = "/sys/devices/platform/ff680020.pwm/pwm/pwmchip2/"
	pwm2pwm0Dir       = pwm2Dir + "pwm0/"
	pwm2ExportPath    = pwm2Dir + "export"
	pwm2UnexportPath  = pwm2Dir + "unexport"
	pwm2EnablePath    = pwm2pwm0Dir + "enable"
	pwm2PeriodPath    = pwm2pwm0Dir + "period"
	pwm2DutyCyclePath = pwm2pwm0Dir + "duty_cycle"
	pwm2PolarityPath  = pwm2pwm0Dir + "polarity"
)

var pwmMockPaths = []string{
	pwm2ExportPath,
	pwm2UnexportPath,
	pwm2EnablePath,
	pwm2PeriodPath,
	pwm2DutyCyclePath,
	pwm2PolarityPath,
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
	fs.Files[pwm2EnablePath].Contents = "0"
	fs.Files[pwm2PeriodPath].Contents = "0"
	fs.Files[pwm2DutyCyclePath].Contents = "0"
	fs.Files[pwm2PolarityPath].Contents = pwmInverted
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

func TestDigitalIO(t *testing.T) {
	// only basic tests needed, further tests are done in "digitalpinsadaptor.go"
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	a.DigitalWrite("7", 1)
	gobottest.Assert(t, fs.Files[gpio17Path+"value"].Contents, "1")

	fs.Files[gpio160Path+"value"].Contents = "1"
	i, _ := a.DigitalRead("10")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("99", 1), errors.New("'99' is not a valid id for a digital pin"))
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

	gobottest.Assert(t, fs.Files[pwm2ExportPath].Contents, "0")
	gobottest.Assert(t, fs.Files[pwm2EnablePath].Contents, "1")
	gobottest.Assert(t, fs.Files[pwm2PeriodPath].Contents, fmt.Sprintf("%d", pwmPeriodDefault))
	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, "3921568")
	gobottest.Assert(t, fs.Files[pwm2PolarityPath].Contents, "normal")

	err = a.ServoWrite("33", 0)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, "500000")

	err = a.ServoWrite("33", 180)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, "2000000")
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
	gobottest.Assert(t, fs.Files[pwm2ExportPath].Contents, "0")
	gobottest.Assert(t, fs.Files[pwm2EnablePath].Contents, "1")
	gobottest.Assert(t, fs.Files[pwm2PeriodPath].Contents, fmt.Sprintf("%d", newPeriod))
	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, "0")
	gobottest.Assert(t, fs.Files[pwm2PolarityPath].Contents, "normal")

	// arrange test for automatic adjustment of duty cycle to lower value
	err = a.PwmWrite("33", 127) // 127 is a little bit smaller than 50% of period
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, fmt.Sprintf("%d", 1270000))
	newPeriod = newPeriod / 10

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, fmt.Sprintf("%d", 127000))

	// arrange test for automatic adjustment of duty cycle to higher value
	newPeriod = newPeriod * 20

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, fmt.Sprintf("%d", 2540000))
}

func TestI2c(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem([]string{"/dev/i2c-1"})
	a.sys.UseMockSyscall()

	con, err := a.GetConnection(0xff, 1)
	gobottest.Assert(t, err, nil)

	_, err = con.Write([]byte{0x00, 0x01})
	gobottest.Assert(t, err, nil)

	data := []byte{42, 42}
	_, err = con.Read(data)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestDefaultBus(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

func TestGetConnectionInvalidBus(t *testing.T) {
	a := NewAdaptor()
	_, err := a.GetConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Bus number 99 out of range"))
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
