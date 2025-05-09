//nolint:nonamedreturns // ok for tests
package adaptors

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/system"
)

const (
	pwmDir             = "/sys/devices/platform/ff680020.pwm/pwm/pwmchip3/" //nolint:gosec // false positive
	pwmPwm44Dir        = pwmDir + "pwm44/"
	pwmPwm47Dir        = pwmDir + "pwm47/"
	pwmExportPath      = pwmDir + "export"
	pwmUnexportPath    = pwmDir + "unexport"
	pwm44EnablePath    = pwmPwm44Dir + "enable"
	pwm44PeriodPath    = pwmPwm44Dir + "period"
	pwm44DutyCyclePath = pwmPwm44Dir + "duty_cycle"
	pwm44PolarityPath  = pwmPwm44Dir + "polarity"
	pwm47EnablePath    = pwmPwm47Dir + "enable"
	pwm47PeriodPath    = pwmPwm47Dir + "period"
	pwm47DutyCyclePath = pwmPwm47Dir + "duty_cycle"
	pwm47PolarityPath  = pwmPwm47Dir + "polarity"
)

var pwmMockPaths = []string{
	pwmExportPath,
	pwmUnexportPath,
	pwm44EnablePath,
	pwm44PeriodPath,
	pwm44DutyCyclePath,
	pwm44PolarityPath,
	pwm47EnablePath,
	pwm47PeriodPath,
	pwm47DutyCyclePath,
	pwm47PolarityPath,
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
	fs.Files[pwm44EnablePath].Contents = "0"
	fs.Files[pwm44PeriodPath].Contents = "0"
	fs.Files[pwm44DutyCyclePath].Contents = "0"
	fs.Files[pwm44PolarityPath].Contents = a.pwmPinsCfg.polarityInvertedIdentifier
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
	channel = channel + 11 // just for tests, 33=>pwm0, 36=>pwm3
	return pwmDir, channel, err
}

func TestNewPWMPinsAdaptor(t *testing.T) {
	// arrange
	translate := func(pin string) (chip string, line int, err error) { return }
	// act
	a := NewPWMPinsAdaptor(system.NewAccesser(), translate)
	// assert
	assert.Equal(t, uint32(pwmPeriodDefault), a.pwmPinsCfg.periodDefault)
	assert.Equal(t, "normal", a.pwmPinsCfg.polarityNormalIdentifier)
	assert.Equal(t, "inversed", a.pwmPinsCfg.polarityInvertedIdentifier)
	assert.True(t, a.pwmPinsCfg.adjustDutyOnSetPeriod)
}

func TestPWMPinsConnect(t *testing.T) {
	translate := func(pin string) (chip string, line int, err error) { return }
	a := NewPWMPinsAdaptor(system.NewAccesser(), translate)
	assert.Equal(t, (map[string]gobot.PWMPinner)(nil), a.pins)

	err := a.PwmWrite("33", 1)
	require.ErrorContains(t, err, "not connected")

	err = a.Connect()
	require.NoError(t, err)
	assert.NotEqual(t, (map[string]gobot.PWMPinner)(nil), a.pins)
	assert.Empty(t, a.pins)
}

func TestPWMPinsFinalize(t *testing.T) {
	// arrange
	sys := system.NewAccesser()
	fs := sys.UseMockFilesystem(pwmMockPaths)
	a := NewPWMPinsAdaptor(sys, testPWMPinTranslator)
	fs.Files[pwm44PeriodPath].Contents = "0"
	fs.Files[pwm44DutyCyclePath].Contents = "0"
	// assert that finalize before connect is working
	require.NoError(t, a.Finalize())
	// arrange
	require.NoError(t, a.Connect())
	require.NoError(t, a.PwmWrite("33", 1))
	assert.Len(t, a.pins, 1)
	// act
	err := a.Finalize()
	// assert
	require.NoError(t, err)
	assert.Empty(t, a.pins)
	// assert that finalize after finalize is working
	require.NoError(t, a.Finalize())
	// arrange missing sysfs file
	require.NoError(t, a.Connect())
	require.NoError(t, a.PwmWrite("33", 2))
	delete(fs.Files, pwmUnexportPath)
	err = a.Finalize()
	require.ErrorContains(t, err, pwmUnexportPath+": no such file")
	// arrange write error
	require.NoError(t, a.Connect())
	require.NoError(t, a.PwmWrite("33", 2))
	fs.WithWriteError = true
	err = a.Finalize()
	require.ErrorContains(t, err, "write error")
}

func TestPWMPinsReConnect(t *testing.T) {
	// arrange
	a, _ := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)
	require.NoError(t, a.PwmWrite("33", 1))
	assert.Len(t, a.pins, 1)
	require.NoError(t, a.Finalize())
	// act
	err := a.Connect()
	// assert
	require.NoError(t, err)
	assert.NotNil(t, a.pins)
	assert.Empty(t, a.pins)
}

func TestPWMPinsCache(t *testing.T) {
	// arrange
	a, _ := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)
	// act
	firstSysPin, err := a.PWMPin("33")
	require.NoError(t, err)
	secondSysPin, err := a.PWMPin("33")
	require.NoError(t, err)
	otherSysPin, err := a.PWMPin("36")
	require.NoError(t, err)
	// assert
	assert.Equal(t, secondSysPin, firstSysPin)
	assert.NotEqual(t, otherSysPin, firstSysPin)
}

func TestPwmWrite(t *testing.T) {
	tests := map[string]struct {
		pin              string
		value            byte
		minimumRate      float64
		simulateWriteErr bool
		simulateReadErr  bool
		wantExport       string
		wantEnable       string
		wantPeriod       string
		wantDutyCycle    string
		wantErr          string
	}{
		"write_max": {
			pin:           "33",
			value:         255,
			wantExport:    "44",
			wantEnable:    "1",
			wantPeriod:    "10000000",
			wantDutyCycle: "10000000",
		},
		"write_nearmax": {
			pin:           "33",
			value:         254,
			wantExport:    "44",
			wantEnable:    "1",
			wantPeriod:    "10000000",
			wantDutyCycle: "9960784",
		},
		"write_mid": {
			pin:           "33",
			value:         100,
			wantExport:    "44",
			wantEnable:    "1",
			wantPeriod:    "10000000",
			wantDutyCycle: "3921568",
		},
		"write_near min": {
			pin:           "33",
			value:         1,
			wantExport:    "44",
			wantEnable:    "1",
			wantPeriod:    "10000000",
			wantDutyCycle: "39215",
		},
		"write_min": {
			pin:           "33",
			value:         0,
			minimumRate:   0.05,
			wantExport:    "44",
			wantEnable:    "1",
			wantPeriod:    "10000000",
			wantDutyCycle: "0",
		},
		"error_min_rate": {
			pin:           "33",
			value:         1,
			minimumRate:   0.05,
			wantExport:    "44",
			wantEnable:    "1",
			wantPeriod:    "10000000",
			wantDutyCycle: "0",
			wantErr:       "is lower than allowed (0.05",
		},
		"error_non_existent_pin": {
			pin:           "notexist",
			wantEnable:    "0",
			wantPeriod:    "0",
			wantDutyCycle: "0",
			wantErr:       "'notexist' is not a valid id of a PWM pin",
		},
		"error_write_error": {
			pin:              "33",
			value:            10,
			simulateWriteErr: true,
			wantEnable:       "0",
			wantPeriod:       "0",
			wantDutyCycle:    "0",
			wantErr:          "write error",
		},
		"error_read_error": {
			pin:             "33",
			value:           11,
			simulateReadErr: true,
			wantExport:      "44",
			wantEnable:      "0",
			wantPeriod:      "0",
			wantDutyCycle:   "0",
			wantErr:         "read error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a, fs := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)
			if tc.minimumRate > 0 {
				a.pwmPinsCfg.dutyRateMinimum = tc.minimumRate
			}
			fs.WithWriteError = tc.simulateWriteErr
			fs.WithReadError = tc.simulateReadErr
			// act
			err := a.PwmWrite(tc.pin, tc.value)
			// assert
			assert.Equal(t, tc.wantExport, fs.Files[pwmExportPath].Contents)
			assert.Equal(t, tc.wantEnable, fs.Files[pwm44EnablePath].Contents)
			assert.Equal(t, tc.wantPeriod, fs.Files[pwm44PeriodPath].Contents)
			assert.Equal(t, tc.wantDutyCycle, fs.Files[pwm44DutyCyclePath].Contents)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "normal", fs.Files[pwm44PolarityPath].Contents)
			}
		})
	}
}

func TestServoWrite(t *testing.T) {
	a, fs := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)

	err := a.ServoWrite("33", 0)
	require.NoError(t, err)
	assert.Equal(t, "44", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwm44EnablePath].Contents)
	//nolint:perfsprint // ok here
	assert.Equal(t, fmt.Sprintf("%d", a.pwmPinsCfg.periodDefault), fs.Files[pwm44PeriodPath].Contents)
	assert.Equal(t, "250000", fs.Files[pwm44DutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwm44PolarityPath].Contents)

	err = a.ServoWrite("33", 180)
	require.NoError(t, err)
	assert.Equal(t, "1250000", fs.Files[pwm44DutyCyclePath].Contents)

	err = a.ServoWrite("notexist", 42)
	require.ErrorContains(t, err, "'notexist' is not a valid id of a PWM pin")

	fs.WithWriteError = true
	err = a.ServoWrite("33", 100)
	require.ErrorContains(t, err, "write error")
	fs.WithWriteError = false

	fs.WithReadError = true
	err = a.ServoWrite("33", 100)
	require.ErrorContains(t, err, "read error")
	fs.WithReadError = false

	delete(a.pwmPinsCfg.pinsServoScale, "33")
	err = a.ServoWrite("33", 42)
	require.EqualError(t, err, "no scaler found for servo pin '33'")
}

func TestSetPeriod(t *testing.T) {
	// arrange
	a, fs := initTestPWMPinsAdaptorWithMockedFilesystem(pwmMockPaths)
	newPeriod := uint32(2550000)
	// act
	err := a.SetPeriod("33", newPeriod)
	// assert
	require.NoError(t, err)
	assert.Equal(t, "44", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwm44EnablePath].Contents)
	assert.Equal(t, fmt.Sprintf("%d", newPeriod), fs.Files[pwm44PeriodPath].Contents) //nolint:perfsprint // ok here
	assert.Equal(t, "0", fs.Files[pwm44DutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwm44PolarityPath].Contents)

	// arrange test for automatic adjustment of duty cycle to lower value
	err = a.PwmWrite("33", 127) // 127 is a little bit smaller than 50% of period
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(1270000), fs.Files[pwm44DutyCyclePath].Contents)
	newPeriod = newPeriod / 10

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(127000), fs.Files[pwm44DutyCyclePath].Contents)

	// arrange test for automatic adjustment of duty cycle to higher value
	newPeriod = newPeriod * 20

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(2540000), fs.Files[pwm44DutyCyclePath].Contents)

	// act
	err = a.SetPeriod("not_exist", newPeriod)
	// assert
	require.ErrorContains(t, err, "'not_exist' is not a valid id of a PWM pin")
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
			mockPaths: []string{pwmExportPath, pwm44EnablePath, pwm44PeriodPath, pwm44DutyCyclePath, pwm44PolarityPath},
			period:    "0",
			dutyCycle: "0",
			translate: translator,
			pin:       "33",
		},
		"init_export_error": {
			mockPaths: []string{},
			translate: translator,
			pin:       "33",
			wantErr: "Export() failed for id 44 with  : " +
				"/sys/devices/platform/ff680020.pwm/pwm/pwmchip3/export: no such file",
		},
		"init_setenabled_error": {
			mockPaths: []string{pwmExportPath, pwm44PeriodPath},
			period:    "1000",
			translate: translator,
			pin:       "33",
			wantErr: "SetEnabled(false) failed for id 44 with  : " +
				"/sys/devices/platform/ff680020.pwm/pwm/pwmchip3/pwm44/enable: no such file",
		},
		"init_setperiod_dutycycle_no_error": {
			mockPaths: []string{pwmExportPath, pwm44EnablePath, pwm44PeriodPath, pwm44DutyCyclePath, pwm44PolarityPath},
			period:    "0",
			dutyCycle: "0",
			translate: translator,
			pin:       "33",
		},
		"init_setperiod_error": {
			mockPaths: []string{pwmExportPath, pwm44EnablePath, pwm44DutyCyclePath},
			dutyCycle: "0",
			translate: translator,
			pin:       "33",
			wantErr: "SetPeriod(10000000) failed for id 44 with  : " +
				"/sys/devices/platform/ff680020.pwm/pwm/pwmchip3/pwm44/period: no such file",
		},
		"init_setpolarity_error": {
			mockPaths: []string{pwmExportPath, pwm44EnablePath, pwm44PeriodPath, pwm44DutyCyclePath},
			period:    "0",
			dutyCycle: "0",
			translate: translator,
			pin:       "33",
			wantErr: "SetPolarity(normal) failed for id 44 with  : " +
				"/sys/devices/platform/ff680020.pwm/pwm/pwmchip3/pwm44/polarity: no such file",
		},
		"translate_error": {
			translate: func(string) (string, int, error) { return "", -1, errors.New(translateErr) },
			wantErr:   translateErr,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			sys := system.NewAccesser()
			fs := sys.UseMockFilesystem(tc.mockPaths)
			if tc.period != "" {
				fs.Files[pwm44PeriodPath].Contents = tc.period
			}
			if tc.dutyCycle != "" {
				fs.Files[pwm44DutyCyclePath].Contents = tc.dutyCycle
			}
			a := NewPWMPinsAdaptor(sys, tc.translate)
			if err := a.Connect(); err != nil {
				panic(err)
			}
			// act
			got, err := a.PWMPin(tc.pin)
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
				assert.NotNil(t, got)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
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
		require.NoError(t, a.Connect())
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
