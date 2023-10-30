package adaptors

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/system"
)

const (
	pwmDir           = "/sys/devices/platform/ff680020.pwm/pwm/pwmchip3/" //nolint:gosec // false positive
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
var (
	_ gobot.PWMPinnerProvider = (*PWMPinsAdaptor)(nil)
	_ gpio.PwmWriter          = (*PWMPinsAdaptor)(nil)
	_ gpio.ServoWriter        = (*PWMPinsAdaptor)(nil)
)

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
	assert.Equal(t, uint32(pwmPeriodDefault), a.periodDefault)
	assert.Equal(t, "normal", a.polarityNormalIdentifier)
	assert.Equal(t, "inverted", a.polarityInvertedIdentifier)
	assert.True(t, a.adjustDutyOnSetPeriod)
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
	assert.Equal(t, wantErr, err)
}

func TestWithPWMPinDefaultPeriod(t *testing.T) {
	// arrange
	const newPeriod = uint32(10)
	a := NewPWMPinsAdaptor(system.NewAccesser(), func(string) (c string, l int, e error) { return })
	// act
	WithPWMPinDefaultPeriod(newPeriod)(a)
	// assert
	assert.Equal(t, newPeriod, a.periodDefault)
}

func TestWithPolarityInvertedIdentifier(t *testing.T) {
	// arrange
	const newPolarityIdent = "pwm_invers"
	a := NewPWMPinsAdaptor(system.NewAccesser(), func(pin string) (c string, l int, e error) { return })
	// act
	WithPolarityInvertedIdentifier(newPolarityIdent)(a)
	// assert
	assert.Equal(t, newPolarityIdent, a.polarityInvertedIdentifier)
}

func TestPWMPinsConnect(t *testing.T) {
	translate := func(pin string) (chip string, line int, err error) { return }
	a := NewPWMPinsAdaptor(system.NewAccesser(), translate)
	assert.Equal(t, (map[string]gobot.PWMPinner)(nil), a.pins)

	err := a.PwmWrite("33", 1)
	assert.ErrorContains(t, err, "not connected")

	err = a.Connect()
	assert.NoError(t, err)
	assert.NotEqual(t, (map[string]gobot.PWMPinner)(nil), a.pins)
	assert.Equal(t, 0, len(a.pins))
}

func TestPWMPinsFinalize(t *testing.T) {
	// arrange
	sys := system.NewAccesser()
	fs := sys.UseMockFilesystem(pwmMockPaths)
	a := NewPWMPinsAdaptor(sys, testPWMPinTranslator)
	fs.Files[pwmPeriodPath].Contents = "0"
	fs.Files[pwmDutyCyclePath].Contents = "0"
	// assert that finalize before connect is working
	assert.NoError(t, a.Finalize())
	// arrange
	assert.NoError(t, a.Connect())
	assert.NoError(t, a.PwmWrite("33", 1))
	assert.Equal(t, 1, len(a.pins))
	// act
	err := a.Finalize()
	// assert
	assert.NoError(t, err)
	assert.Equal(t, 0, len(a.pins))
	// assert that finalize after finalize is working
	assert.NoError(t, a.Finalize())
	// arrange missing sysfs file
	assert.NoError(t, a.Connect())
	assert.NoError(t, a.PwmWrite("33", 2))
	delete(fs.Files, pwmUnexportPath)
	err = a.Finalize()
	assert.Contains(t, err.Error(), pwmUnexportPath+": no such file")
	// arrange write error
	assert.NoError(t, a.Connect())
	assert.NoError(t, a.PwmWrite("33", 2))
	fs.WithWriteError = true
	err = a.Finalize()
	assert.Contains(t, err.Error(), "write error")
}

func TestPWMPinsReConnect(t *testing.T) {
	// arrange
	a, _ := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)
	assert.NoError(t, a.PwmWrite("33", 1))
	assert.Equal(t, 1, len(a.pins))
	assert.NoError(t, a.Finalize())
	// act
	err := a.Connect()
	// assert
	assert.NoError(t, err)
	assert.NotNil(t, a.pins)
	assert.Equal(t, 0, len(a.pins))
}

func TestPwmWrite(t *testing.T) {
	a, fs := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)

	err := a.PwmWrite("33", 100)
	assert.NoError(t, err)

	assert.Equal(t, "44", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwmEnablePath].Contents)
	assert.Equal(t, fmt.Sprintf("%d", a.periodDefault), fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "3921568", fs.Files[pwmDutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwmPolarityPath].Contents)

	err = a.PwmWrite("notexist", 42)
	assert.ErrorContains(t, err, "'notexist' is not a valid id of a PWM pin")

	fs.WithWriteError = true
	err = a.PwmWrite("33", 100)
	assert.Contains(t, err.Error(), "write error")
	fs.WithWriteError = false

	fs.WithReadError = true
	err = a.PwmWrite("33", 100)
	assert.Contains(t, err.Error(), "read error")
}

func TestServoWrite(t *testing.T) {
	a, fs := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)

	err := a.ServoWrite("33", 0)
	assert.Equal(t, "44", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwmEnablePath].Contents)
	assert.Equal(t, fmt.Sprintf("%d", a.periodDefault), fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "500000", fs.Files[pwmDutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwmPolarityPath].Contents)
	assert.NoError(t, err)

	err = a.ServoWrite("33", 180)
	assert.NoError(t, err)
	assert.Equal(t, "2000000", fs.Files[pwmDutyCyclePath].Contents)

	err = a.ServoWrite("notexist", 42)
	assert.ErrorContains(t, err, "'notexist' is not a valid id of a PWM pin")

	fs.WithWriteError = true
	err = a.ServoWrite("33", 100)
	assert.Contains(t, err.Error(), "write error")
	fs.WithWriteError = false

	fs.WithReadError = true
	err = a.ServoWrite("33", 100)
	assert.Contains(t, err.Error(), "read error")
}

func TestSetPeriod(t *testing.T) {
	// arrange
	a, fs := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)
	newPeriod := uint32(2550000)
	// act
	err := a.SetPeriod("33", newPeriod)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, "44", fs.Files[pwmExportPath].Contents)
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

	// act
	err = a.SetPeriod("not_exist", newPeriod)
	// assert
	assert.ErrorContains(t, err, "'not_exist' is not a valid id of a PWM pin")
}

func Test_PWMPin(t *testing.T) {
	translateErr := "translator_error"
	translator := func(string) (string, int, error) { return pwmDir, 44, nil }
	tests := map[string]struct {
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
				assert.NoError(t, err)
				assert.NotNil(t, got)
			} else {
				if !strings.Contains(err.Error(), tc.wantErr) {
					log.Println(err.Error())
				}
				assert.Contains(t, err.Error(), tc.wantErr)
				assert.Nil(t, got)
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
