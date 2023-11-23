//nolint:forcetypeassert // ok here
package aio

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemperatureSensorDriver(t *testing.T) {
	// arrange
	const pin = "123"
	a := newAioTestAdaptor()
	// act
	d := NewTemperatureSensorDriver(a, pin)
	// assert: driver attributes
	assert.IsType(t, &TemperatureSensorDriver{}, d)
	assert.NotNil(t, d.driverCfg)
	assert.True(t, strings.HasPrefix(d.Name(), "TemperatureSensor"))
	assert.Equal(t, a, d.Connection())
	require.NoError(t, d.afterStart())
	require.NoError(t, d.beforeHalt())
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: sensor attributes
	assert.Equal(t, pin, d.Pin())
	assert.InDelta(t, 0.0, d.lastValue, 0, 0)
	assert.Equal(t, 0, d.lastRawValue)
	assert.NotNil(t, d.halt)
	assert.NotNil(t, d.Eventer)
	require.NotNil(t, d.sensorCfg)
	assert.Equal(t, time.Duration(0), d.sensorCfg.readInterval)
	assert.NotNil(t, d.sensorCfg.scale)
}

func TestNewTemperatureSensorDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName     = "outlet temperature"
		cycReadDur = 10 * time.Millisecond
	)
	panicFunc := func() {
		NewTemperatureSensorDriver(newAioTestAdaptor(), "1", WithName("crazy"),
			WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewTemperatureSensorDriver(newAioTestAdaptor(), "1", WithName(myName), WithSensorCyclicRead(cycReadDur))
	// assert
	assert.Equal(t, cycReadDur, d.sensorCfg.readInterval)
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestTemperatureSensorRead_NtcScaler(t *testing.T) {
	tests := map[string]struct {
		input int
		want  float64
	}{
		"smaller_than_min": {input: -1, want: 457.720219684306},
		"min":              {input: 0, want: 457.720219684306},
		"near_min":         {input: 1, want: 457.18923673420545},
		"mid_range":        {input: 127, want: 87.9784401845593},
		"T25C":             {input: 232, want: 24.805280460718336},
		"T0C":              {input: 248, want: -0.9858175109026774},
		"T-25C":            {input: 253, want: -22.92863536929451},
		"near_max":         {input: 254, want: -33.51081663999781},
		"max":              {input: 255, want: -273.15},
		"bigger_than_max":  {input: 256, want: -273.15},
	}
	a := newAioTestAdaptor()
	d := NewTemperatureSensorDriver(a, "4")
	ntc1 := TemperatureSensorNtcConf{TC0: 25, R0: 10000.0, B: 3950} // Ohm, R25=10k, B=3950
	d.SetNtcScaler(255, 1000, true, ntc1)                           // Ohm, reference value: 3300, series R: 1k
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a.analogReadFunc = func() (int, error) {
				return tc.input, nil
			}
			// act
			got, err := d.Read()
			// assert
			require.NoError(t, err)
			assert.InDelta(t, tc.want, got, 0.0)
		})
	}
}

func TestTemperatureSensorDriver_LinearScaler(t *testing.T) {
	tests := map[string]struct {
		input int
		want  float64
	}{
		"smaller_than_min": {input: -129, want: -40},
		"min":              {input: -128, want: -40},
		"near_min":         {input: -127, want: -39.450980392156865},
		"T-25C":            {input: -101, want: -25.17647058823529},
		"T0C":              {input: -55, want: 0.07843137254902288},
		"T25C":             {input: -10, want: 24.7843137254902},
		"mid_range":        {input: 0, want: 30.274509803921575},
		"near_max":         {input: 126, want: 99.45098039215688},
		"max":              {input: 127, want: 100},
		"bigger_than_max":  {input: 128, want: 100},
	}
	a := newAioTestAdaptor()
	d := NewTemperatureSensorDriver(a, "4")
	d.SetLinearScaler(-128, 127, -40, 100)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a.analogReadFunc = func() (int, error) {
				return tc.input, nil
			}
			// act
			got, err := d.Read()
			// assert
			require.NoError(t, err)
			assert.InDelta(t, tc.want, got, 0.0)
		})
	}
}

func TestTemperatureSensorPublishesTemperatureInCelsius(t *testing.T) {
	// arrange
	sem := make(chan bool)
	a := newAioTestAdaptor()
	d := NewTemperatureSensorDriver(a, "1", WithSensorCyclicRead(10*time.Millisecond))
	ntc := TemperatureSensorNtcConf{TC0: 25, R0: 10000.0, B: 3975} // Ohm, R25=10k
	d.SetNtcScaler(1023, 10000, false, ntc)                        // Ohm, reference value: 1023, series R: 10k

	a.analogReadFunc = func() (int, error) {
		return 585, nil
	}
	_ = d.Once(d.Event(Value), func(data interface{}) {
		assert.Equal(t, "31.62", fmt.Sprintf("%.2f", data.(float64)))
		sem <- true
	})
	require.NoError(t, d.Start())

	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf(" Temperature Sensor Event \"Data\" was not published")
	}

	assert.InDelta(t, 31.61532462352477, d.Value(), 0.0)
}

func TestTemperatureSensorWithSensorCyclicReadPublishesError(t *testing.T) {
	// arrange
	sem := make(chan bool)
	a := newAioTestAdaptor()
	d := NewTemperatureSensorDriver(a, "1", WithSensorCyclicRead(10*time.Millisecond))

	// send error
	a.analogReadFunc = func() (int, error) {
		return 0, errors.New("read error")
	}

	require.NoError(t, d.Start())

	// expect error
	_ = d.Once(d.Event(Error), func(data interface{}) {
		assert.Equal(t, "read error", data.(error).Error())
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf(" Temperature Sensor Event \"Error\" was not published")
	}
}

func TestTemperatureSensorHalt_WithSensorCyclicRead(t *testing.T) {
	// arrange
	d := NewTemperatureSensorDriver(newAioTestAdaptor(), "1", WithSensorCyclicRead(10*time.Millisecond))
	done := make(chan struct{})
	go func() {
		<-d.halt
		close(done)
	}()
	// act & assert
	require.NoError(t, d.Halt())
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Errorf(" Temperature Sensor was not halted")
	}
}

func TestTempDriver_initialize(t *testing.T) {
	tests := map[string]struct {
		input TemperatureSensorNtcConf
		want  TemperatureSensorNtcConf
	}{
		"B_low_tc0": {
			input: TemperatureSensorNtcConf{TC0: -13, B: 2601.5},
			want:  TemperatureSensorNtcConf{TC0: -13, B: 2601.5, t0: 260.15, r: 10},
		},
		"B_low_tc0_B": {
			input: TemperatureSensorNtcConf{TC0: -13, B: 5203},
			want:  TemperatureSensorNtcConf{TC0: -13, B: 5203, t0: 260.15, r: 20},
		},
		"B_mid_tc0": {
			input: TemperatureSensorNtcConf{TC0: 25, B: 3950},
			want:  TemperatureSensorNtcConf{TC0: 25, B: 3950, t0: 298.15, r: 13.248364916988095},
		},
		"B_mid_tc0_r0_no_change": {
			input: TemperatureSensorNtcConf{TC0: 25, R0: 1234.5, B: 3950},
			want:  TemperatureSensorNtcConf{TC0: 25, R0: 1234.5, B: 3950, t0: 298.15, r: 13.248364916988095},
		},
		"B_high_tc0": {
			input: TemperatureSensorNtcConf{TC0: 100, B: 3731.5},
			want:  TemperatureSensorNtcConf{TC0: 100, B: 3731.5, t0: 373.15, r: 10},
		},
		"T1_low": {
			input: TemperatureSensorNtcConf{TC0: 25, R0: 2500.0, TC1: -13, R1: 10000},
			want: TemperatureSensorNtcConf{
				TC0: 25, R0: 2500.0, TC1: -13, R1: 10000, B: 2829.6355560320544, t0: 298.15, r: 9.490644159087891,
			},
		},
		"T1_high": {
			input: TemperatureSensorNtcConf{TC0: 25, R0: 2500.0, TC1: 100, R1: 371},
			want: TemperatureSensorNtcConf{
				TC0: 25, R0: 2500.0, TC1: 100, R1: 371, B: 2830.087381913779, t0: 298.15, r: 9.49215959052081,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			ntc := tc.input
			// act
			ntc.initialize()
			// assert
			assert.Equal(t, tc.want, ntc)
		})
	}
}
