package adaptors

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/gobottest"
	"gobot.io/x/gobot/v2/system"
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
	fs.Files[pwmPolarityPath].Contents = a.polarityInvertedIdentifier
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

func TestNewPWMPinsAdaptor(t *testing.T) {
	// arrange
	translate := func(pin string) (chip string, line int, err error) { return }
	// act
	a := NewPWMPinsAdaptor(system.NewAccesser(), translate)
	// assert
	gobottest.Assert(t, a.periodDefault, uint32(pwmPeriodDefault))
	gobottest.Assert(t, a.polarityNormalIdentifier, "normal")
	gobottest.Assert(t, a.polarityInvertedIdentifier, "inverted")
	gobottest.Assert(t, a.adjustDutyOnSetPeriod, true)
}

func TestWithPWMPinInitializer(t *testing.T) {
	// This is a general test, that options are applied by using the WithPWMPinInitializer() option.
	// All other configuration options can also be tested by With..(val)(a).
	// arrange
	wantErr := fmt.Errorf("new_initializer")
	newInitializer := func(gobot.PWMPinner) error { return wantErr }
	// act
	a := NewPWMPinsAdaptor(system.NewAccesser(), func(pin string) (c string, l int, e error) { return },
		WithPWMPinInitializer(newInitializer))
	// assert
	err := a.initialize(nil)
	gobottest.Assert(t, err, wantErr)
}

func TestWithPWMPinDefaultPeriod(t *testing.T) {
	// arrange
	const newPeriod = uint32(10)
	a := NewPWMPinsAdaptor(system.NewAccesser(), func(string) (c string, l int, e error) { return })
	// act
	WithPWMPinDefaultPeriod(newPeriod)(a)
	// assert
	gobottest.Assert(t, a.periodDefault, newPeriod)
}

func TestWithPolarityInvertedIdentifier(t *testing.T) {
	// arrange
	const newPolarityIdent = "pwm_invers"
	a := NewPWMPinsAdaptor(system.NewAccesser(), func(pin string) (c string, l int, e error) { return })
	// act
	WithPolarityInvertedIdentifier(newPolarityIdent)(a)
	// assert
	gobottest.Assert(t, a.polarityInvertedIdentifier, newPolarityIdent)
}

func TestPWMPinsConnect(t *testing.T) {
	translate := func(pin string) (chip string, line int, err error) { return }
	a := NewPWMPinsAdaptor(system.NewAccesser(), translate)
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
	fs.Files[pwmPeriodPath].Contents = "0"
	fs.Files[pwmDutyCyclePath].Contents = "0"
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
	gobottest.Assert(t, strings.Contains(err.Error(), pwmUnexportPath+": no such file"), true)
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
	gobottest.Refute(t, a.pins, nil)
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

func TestServoWrite(t *testing.T) {
	a, fs := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)

	err := a.ServoWrite("33", 0)
	gobottest.Assert(t, fs.Files[pwmExportPath].Contents, "44")
	gobottest.Assert(t, fs.Files[pwmEnablePath].Contents, "1")
	gobottest.Assert(t, fs.Files[pwmPeriodPath].Contents, fmt.Sprintf("%d", a.periodDefault))
	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, "500000")
	gobottest.Assert(t, fs.Files[pwmPolarityPath].Contents, "normal")
	gobottest.Assert(t, err, nil)

	err = a.ServoWrite("33", 180)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwmDutyCyclePath].Contents, "2000000")

	err = a.ServoWrite("notexist", 42)
	gobottest.Assert(t, err.Error(), "'notexist' is not a valid id of a PWM pin")

	fs.WithWriteError = true
	err = a.ServoWrite("33", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
	fs.WithWriteError = false

	fs.WithReadError = true
	err = a.ServoWrite("33", 100)
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

	// act
	err = a.SetPeriod("not_exist", newPeriod)
	// assert
	gobottest.Assert(t, err.Error(), "'not_exist' is not a valid id of a PWM pin")
}

func Test_PWMPin(t *testing.T) {
	translateErr := "translator_error"
	translator := func(string) (string, int, error) { return pwmDir, 44, nil }
	var tests = map[string]struct {
		mockPaths []string
		period    string
		dutyCycle string
		translate func(string) (string, int, error)
		pin       string
		wantErr   string
	}{
		"pin_ok": {
			mockPaths: []string{pwmExportPath, pwmEnablePath, pwmPeriodPath, pwmDutyCyclePath, pwmPolarityPath},
			period:    "0",
			dutyCycle: "0",
			translate: translator,
			pin:       "33",
		},
		"init_export_error": {
			mockPaths: []string{},
			translate: translator,
			pin:       "33",
			wantErr:   "Export() failed for id 44 with  : /sys/devices/platform/ff680020.pwm/pwm/pwmchip3/export: no such file",
		},
		"init_setenabled_error": {
			mockPaths: []string{pwmExportPath, pwmPeriodPath},
			period:    "1000",
			translate: translator,
			pin:       "33",
			wantErr:   "SetEnabled(false) failed for id 44 with  : /sys/devices/platform/ff680020.pwm/pwm/pwmchip3/pwm44/enable: no such file",
		},
		"init_setperiod_dutycycle_no_error": {
			mockPaths: []string{pwmExportPath, pwmEnablePath, pwmPeriodPath, pwmDutyCyclePath, pwmPolarityPath},
			period:    "0",
			dutyCycle: "0",
			translate: translator,
			pin:       "33",
		},
		"init_setperiod_error": {
			mockPaths: []string{pwmExportPath, pwmEnablePath, pwmDutyCyclePath},
			dutyCycle: "0",
			translate: translator,
			pin:       "33",
			wantErr:   "SetPeriod(10000000) failed for id 44 with  : /sys/devices/platform/ff680020.pwm/pwm/pwmchip3/pwm44/period: no such file",
		},
		"init_setpolarity_error": {
			mockPaths: []string{pwmExportPath, pwmEnablePath, pwmPeriodPath, pwmDutyCyclePath},
			period:    "0",
			dutyCycle: "0",
			translate: translator,
			pin:       "33",
			wantErr:   "SetPolarity(normal) failed for id 44 with  : /sys/devices/platform/ff680020.pwm/pwm/pwmchip3/pwm44/polarity: no such file",
		},
		"translate_error": {
			translate: func(string) (string, int, error) { return "", -1, fmt.Errorf(translateErr) },
			wantErr:   translateErr,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			sys := system.NewAccesser()
			fs := sys.UseMockFilesystem(tc.mockPaths)
			if tc.period != "" {
				fs.Files[pwmPeriodPath].Contents = tc.period
			}
			if tc.dutyCycle != "" {
				fs.Files[pwmDutyCyclePath].Contents = tc.dutyCycle
			}
			a := NewPWMPinsAdaptor(sys, tc.translate)
			if err := a.Connect(); err != nil {
				panic(err)
			}
			// act
			got, err := a.PWMPin(tc.pin)
			// assert
			if tc.wantErr == "" {
				gobottest.Assert(t, err, nil)
				gobottest.Refute(t, got, nil)
			} else {
				if !strings.Contains(err.Error(), tc.wantErr) {
					log.Println(err.Error())
				}
				gobottest.Assert(t, strings.Contains(err.Error(), tc.wantErr), true)
				gobottest.Assert(t, got, nil)
			}

		})
	}
}

func TestPWMPinConcurrency(t *testing.T) {
	oldProcs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(8)
	defer runtime.GOMAXPROCS(oldProcs)

	translate := func(pin string) (string, int, error) { line, err := strconv.Atoi(pin); return "", line, err }
	sys := system.NewAccesser()

	for retry := 0; retry < 20; retry++ {

		a := NewPWMPinsAdaptor(sys, translate)
		_ = a.Connect()
		var wg sync.WaitGroup

		for i := 0; i < 20; i++ {
			wg.Add(1)
			pinAsString := strconv.Itoa(i)
			go func(pin string) {
				defer wg.Done()
				_, _ = a.PWMPin(pin)
			}(pinAsString)
		}

		wg.Wait()
	}
}
