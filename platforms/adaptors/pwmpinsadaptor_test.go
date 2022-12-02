package adaptors

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/system"
)

const (
	pwmDir           = "/sys/devices/platform/ff680020.pwm/pwm/pwmchip3/"
	pwmPwm0Dir       = pwmDir + "pwm44/"
	pwmExportPath    = pwmDir + "export"
	pwmUnexportPath  = pwmDir + "unexport"
	pwmEnablePath    = pwmPwm0Dir + "enable"
	pwmPeriodPath    = pwmPwm0Dir + "period"
	pwmDutyCyclePath = pwmPwm0Dir + "duty_cycle"
	pwmPolarityPath  = pwmPwm0Dir + "polarity"
)

var pwmMockPaths = []string{
	pwmExportPath,
	pwmUnexportPath,
	pwmEnablePath,
	pwmPeriodPath,
	pwmDutyCyclePath,
	pwmPolarityPath,
}

// make sure that this PWMPinsAdaptor fulfills all the required interfaces
var _ gobot.PWMPinnerProvider = (*PWMPinsAdaptor)(nil)
var _ gpio.PwmWriter = (*PWMPinsAdaptor)(nil)
var _ gpio.ServoWriter = (*PWMPinsAdaptor)(nil)

func initTestPWMPinsAdaptorWithMockedFilesystem(mockPaths []string) (*PWMPinsAdaptor, *system.MockFilesystem) {
	sys := system.NewAccesser()
	fs := sys.UseMockFilesystem(mockPaths)
	a := NewPWMPinsAdaptor(sys, testPWMPinTranslator)
	fs.Files[pwmEnablePath].Contents = "0"
	fs.Files[pwmPeriodPath].Contents = "0"
	fs.Files[pwmDutyCyclePath].Contents = "0"
	fs.Files[pwmPolarityPath].Contents = a.polarityInverted
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func testPWMPinTranslator(id string) (string, int, error) {
	channel, err := strconv.Atoi(id)
	if err != nil {
		return "", -1, fmt.Errorf("'%s' is not a valid id of a PWM pin", id)
	}
	channel = channel + 11 // just for tests
	return pwmDir, channel, err
}

// TODO: test With...

func TestPWMPinsConnect(t *testing.T) {
	translate := func(pin string) (chip string, line int, err error) { return }
	sys := system.NewAccesser()

	a := NewPWMPinsAdaptor(sys, translate)
	gobottest.Assert(t, a.pins, (map[string]gobot.PWMPinner)(nil))

	err := a.PwmWrite("33", 1)
	gobottest.Assert(t, err.Error(), "not connected")

	err = a.Connect()
	gobottest.Assert(t, err, nil)
	gobottest.Refute(t, a.pins, (map[string]gobot.PWMPinner)(nil))
	gobottest.Assert(t, len(a.pins), 0)
}

func TestPWMPinsFinalize(t *testing.T) {
	// arrange
	sys := system.NewAccesser()
	fs := sys.UseMockFilesystem(pwmMockPaths)
	a := NewPWMPinsAdaptor(sys, testPWMPinTranslator)
	// assert that finalize before connect is working
	gobottest.Assert(t, a.Finalize(), nil)
	// arrange
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.PwmWrite("33", 1), nil)
	gobottest.Assert(t, len(a.pins), 1)
	// act
	err := a.Finalize()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pins), 0)
	// assert that finalize after finalize is working
	gobottest.Assert(t, a.Finalize(), nil)
	// arrange missing sysfs file
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.PwmWrite("33", 2), nil)
	delete(fs.Files, pwmUnexportPath)
	err = a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), pwmUnexportPath+": No such file"), true)
	// arrange write error
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.PwmWrite("33", 2), nil)
	fs.WithWriteError = true
	err = a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestPWMPinsReConnect(t *testing.T) {
	// arrange
	a, _ := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)
	gobottest.Assert(t, a.PwmWrite("33", 1), nil)
	gobottest.Assert(t, len(a.pins), 1)
	gobottest.Assert(t, a.Finalize(), nil)
	// act
	err := a.Connect()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pins), 0)
}

func TestPwmWrite(t *testing.T) {
	a, fs := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)

	err := a.PwmWrite("33", 100)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files[pwmExportPath].Contents, "44")
	gobottest.Assert(t, fs.Files[pwmEnablePath].Contents, "1")
	gobottest.Assert(t, fs.Files[pwmPeriodPath].Contents, fmt.Sprintf("%d", a.periodDefault))
	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, "3921568")
	gobottest.Assert(t, fs.Files[pwmPolarityPath].Contents, "normal")

	err = a.ServoWrite("33", 0)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, "500000")

	err = a.ServoWrite("33", 180)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, "2000000")

	err = a.PwmWrite("notexist", 42)
	gobottest.Assert(t, err.Error(), "'notexist' is not a valid id of a PWM pin")

	fs.WithWriteError = true
	err = a.PwmWrite("33", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
	fs.WithWriteError = false

	fs.WithReadError = true
	err = a.PwmWrite("33", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "read error"), true)
}

func TestSetPeriod(t *testing.T) {
	// arrange
	a, fs := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)
	newPeriod := uint32(2550000)
	// act
	err := a.SetPeriod("33", newPeriod)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwmExportPath].Contents, "44")
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
