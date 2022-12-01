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

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*PWMPinsAdaptor, *system.MockFilesystem) {
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

func TestInvalidPWMPin(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem(pwmMockPaths)

	err := a.PwmWrite("notexist", 42)
	gobottest.Assert(t, err.Error(), "'notexist' is not a valid id of a PWM pin")
}

func TestPwmWrite(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)

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
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestPwmWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	fs.WithWriteError = true

	err := a.PwmWrite("33", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestPwmWriteReadError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	fs.WithReadError = true

	err := a.PwmWrite("33", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "read error"), true)
}

func TestSetPeriod(t *testing.T) {
	// arrange
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
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

func TestFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)

	gobottest.Assert(t, a.PwmWrite("33", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}
