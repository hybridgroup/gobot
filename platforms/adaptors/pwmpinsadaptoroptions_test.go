//nolint:nonamedreturns // ok for tests
package adaptors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

func TestWithPWMPinInitializer(t *testing.T) {
	// This is a general test, that options are applied by using the WithPWMPinInitializer() option.
	// All other configuration options can also be tested by With..(val).apply(cfg).
	// arrange
	wantErr := fmt.Errorf("new_initializer")
	newInitializer := func(string, gobot.PWMPinner) error { return wantErr }
	// act
	a := NewPWMPinsAdaptor(system.NewAccesser(), func(pin string) (c string, l int, e error) { return },
		WithPWMPinInitializer(newInitializer))
	// assert
	err := a.pwmPinsCfg.initialize("1", nil)
	assert.Equal(t, wantErr, err)
}

func TestWithPWMDefaultPeriod(t *testing.T) {
	// arrange
	const newPeriod = uint32(10)
	cfg := &pwmPinsConfiguration{periodDefault: 123}
	// act
	WithPWMDefaultPeriod(newPeriod).apply(cfg)
	// assert
	assert.Equal(t, newPeriod, cfg.periodDefault)
}

func TestWithPWMPolarityInvertedIdentifier(t *testing.T) {
	// arrange
	const newPolarityIdent = "pwm_invers"
	cfg := &pwmPinsConfiguration{polarityInvertedIdentifier: "old_inverted"}
	// act
	WithPWMPolarityInvertedIdentifier(newPolarityIdent).apply(cfg)
	// assert
	assert.Equal(t, newPolarityIdent, cfg.polarityInvertedIdentifier)
}

func TestWithPWMNoDutyCycleAdjustment(t *testing.T) {
	// arrange
	cfg := &pwmPinsConfiguration{adjustDutyOnSetPeriod: true}
	// act
	WithPWMNoDutyCycleAdjustment().apply(cfg)
	// assert
	assert.False(t, cfg.adjustDutyOnSetPeriod)
}

func TestWithPWMDefaultPeriodForPin(t *testing.T) {
	// arrange
	const (
		pin       = "pin4test"
		newPeriod = 123456
	)
	cfg := &pwmPinsConfiguration{pinsDefaultPeriod: map[string]uint32{pin: 54321}}
	// act
	WithPWMDefaultPeriodForPin(pin, newPeriod).apply(cfg)
	// assert
	assert.Equal(t, uint32(newPeriod), cfg.pinsDefaultPeriod[pin])
}

func TestWithPWMServoDutyCycleRangeForPin(t *testing.T) {
	const (
		pin    = "pin4test"
		newMin = 19
		newMax = 99
	)

	tests := map[string]struct {
		scaleMap     map[string]pwmPinServoScale
		wantScaleMap map[string]pwmPinServoScale
	}{
		"empty_scale_map": {
			scaleMap: make(map[string]pwmPinServoScale),
			wantScaleMap: map[string]pwmPinServoScale{
				pin: {minDuty: newMin, maxDuty: newMax, minDegree: 0, maxDegree: 180},
			},
		},
		"scale_exists_for_set_pin": {
			scaleMap: map[string]pwmPinServoScale{
				"other": {minDuty: 123, maxDuty: 234, minDegree: 11, maxDegree: 22},
				pin:     {minDuty: newMin - 2, maxDuty: newMax + 2, minDegree: 1, maxDegree: 2},
			},
			wantScaleMap: map[string]pwmPinServoScale{
				"other": {minDuty: 123, maxDuty: 234, minDegree: 11, maxDegree: 22},
				pin:     {minDuty: newMin, maxDuty: newMax, minDegree: 1, maxDegree: 2},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			cfg := &pwmPinsConfiguration{pinsServoScale: tc.scaleMap}
			// act
			WithPWMServoDutyCycleRangeForPin(pin, newMin, newMax).apply(cfg)
			// assert
			assert.Equal(t, tc.wantScaleMap, cfg.pinsServoScale)
		})
	}
}

func TestWithPWMServoAngleRangeForPin(t *testing.T) {
	const (
		pin    = "pin4test"
		newMin = 30
		newMax = 90
	)

	tests := map[string]struct {
		scaleMap     map[string]pwmPinServoScale
		wantScaleMap map[string]pwmPinServoScale
	}{
		"empty_scale_map": {
			scaleMap: make(map[string]pwmPinServoScale),
			wantScaleMap: map[string]pwmPinServoScale{
				pin: {minDuty: 0.0, maxDuty: 0.0, minDegree: newMin, maxDegree: newMax},
			},
		},
		"scale_exists_for_set_pin": {
			scaleMap: map[string]pwmPinServoScale{
				"other": {minDuty: 123, maxDuty: 234, minDegree: 11, maxDegree: 22},
				pin:     {minDuty: 4, maxDuty: 5, minDegree: newMin - 2, maxDegree: newMax + 2},
			},
			wantScaleMap: map[string]pwmPinServoScale{
				"other": {minDuty: 123, maxDuty: 234, minDegree: 11, maxDegree: 22},
				pin:     {minDuty: 4, maxDuty: 5, minDegree: newMin, maxDegree: newMax},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			cfg := &pwmPinsConfiguration{pinsServoScale: tc.scaleMap}
			// act
			WithPWMServoAngleRangeForPin(pin, newMin, newMax).apply(cfg)
			// assert
			assert.Equal(t, tc.wantScaleMap, cfg.pinsServoScale)
		})
	}
}
