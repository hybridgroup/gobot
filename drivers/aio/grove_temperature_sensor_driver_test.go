//nolint:forcetypeassert // ok here
package aio

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*GroveTemperatureSensorDriver)(nil)

func TestNewGroveTemperatureSensorDriver(t *testing.T) {
	// arrange
	const pin = "123"
	a := newAioTestAdaptor()
	// act
	d := NewGroveTemperatureSensorDriver(a, pin)
	// assert: driver attributes
	assert.IsType(t, &GroveTemperatureSensorDriver{}, d)
	assert.NotNil(t, d.driverCfg)
	assert.True(t, strings.HasPrefix(d.Name(), "GroveTemperatureSensor"))
	assert.Equal(t, a, d.Connection())
	require.NoError(t, d.afterStart())
	require.NoError(t, d.beforeHalt())
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: sensor attributes
	assert.Equal(t, pin, d.Pin())
	assert.InDelta(t, 0.0, d.lastValue, 0, 0)
	assert.Equal(t, 0, d.lastRawValue)
	assert.Nil(t, d.halt) // will be created on initialize, if cyclic reading is on
	assert.NotNil(t, d.Eventer)
	require.NotNil(t, d.sensorCfg)
	assert.Equal(t, time.Duration(0), d.sensorCfg.readInterval)
	assert.NotNil(t, d.sensorCfg.scale)
}

func TestNewGroveTemperatureSensorDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName     = "inlet temperature"
		cycReadDur = 10 * time.Millisecond
	)
	panicFunc := func() {
		NewGroveTemperatureSensorDriver(newAioTestAdaptor(), "1", WithName("crazy"),
			WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewGroveTemperatureSensorDriver(newAioTestAdaptor(), "1", WithName(myName), WithSensorCyclicRead(cycReadDur))
	// assert
	assert.Equal(t, cycReadDur, d.sensorCfg.readInterval)
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestGroveTemperatureSensorRead_scaler(t *testing.T) {
	tests := map[string]struct {
		input int
		want  float64
	}{
		"min":           {input: 0, want: -273.15},
		"nearMin":       {input: 1, want: -76.96736464322436},
		"T-25C":         {input: 65, want: -25.064097201780044},
		"T0C":           {input: 233, want: -0.014379114122164083},
		"T25C":          {input: 511, want: 24.956285721537938},
		"585":           {input: 585, want: 31.61532462352477},
		"nearMax":       {input: 1022, want: 347.6819764792606},
		"max":           {input: 1023, want: 347.77682140097613},
		"biggerThanMax": {input: 5000, want: 347.77682140097613},
	}
	a := newAioTestAdaptor()
	d := NewGroveTemperatureSensorDriver(a, "54")
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

func TestGroveTemperatureSensor_publishesTemperatureInCelsius(t *testing.T) {
	// arrange
	sem := make(chan bool)
	a := newAioTestAdaptor()
	d := NewGroveTemperatureSensorDriver(a, "1", WithSensorCyclicRead(10*time.Millisecond))

	a.analogReadFunc = func() (int, error) {
		return 585, nil
	}

	// act: start cyclic reading
	require.NoError(t, d.Start())

	_ = d.Once(d.Event(Value), func(data interface{}) {
		assert.Equal(t, "31.62", fmt.Sprintf("%.2f", data.(float64)))
		sem <- true
	})

	// assert: value was published
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		require.Fail(t, "Grove Temperature Sensor Event \"Value\" was not published")
	}

	assert.InDelta(t, 31.61532462352477, d.Temperature(), 0.0)
}
