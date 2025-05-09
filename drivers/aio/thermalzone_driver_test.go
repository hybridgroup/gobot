package aio

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewThermalZoneDriver(t *testing.T) {
	// arrange
	const pin = "thermal_zone0"
	a := newAioTestAdaptor()
	// act
	d := NewThermalZoneDriver(a, pin)
	// assert: driver attributes
	assert.IsType(t, &ThermalZoneDriver{}, d)
	assert.NotNil(t, d.driverCfg)
	assert.True(t, strings.HasPrefix(d.Name(), "ThermalZone"))
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
	// assert: thermal zone attributes
	require.NotNil(t, d.thermalZoneCfg)
	require.NotNil(t, d.thermalZoneCfg.scaleUnit)
	assert.InDelta(t, 1.0, d.thermalZoneCfg.scaleUnit(1), 0.0)
}

func TestNewThermalZoneDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName     = "outlet temperature"
		cycReadDur = 10 * time.Millisecond
	)
	panicFunc := func() {
		NewThermalZoneDriver(newAioTestAdaptor(), "1", WithName("crazy"),
			WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewThermalZoneDriver(newAioTestAdaptor(), "1",
		WithName(myName),
		WithSensorCyclicRead(cycReadDur),
		WithFahrenheit())
	// assert
	assert.Equal(t, cycReadDur, d.sensorCfg.readInterval)
	assert.InDelta(t, 33.8, d.thermalZoneCfg.scaleUnit(1), 0.0) // (1°C × 9/5) + 32 = 33,8°F
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestThermalZoneWithSensorCyclicRead_PublishesTemperatureInFahrenheit(t *testing.T) {
	// arrange
	sem := make(chan bool)
	a := newAioTestAdaptor()
	d := NewThermalZoneDriver(a, "1", WithSensorCyclicRead(10*time.Millisecond), WithFahrenheit())
	a.analogReadFunc = func() (int, error) {
		return -100000, nil // -100.000 °C
	}
	// act: start cyclic reading
	require.NoError(t, d.Start())
	// assert
	_ = d.Once(d.Event(Value), func(data interface{}) {
		//nolint:forcetypeassert // ok here
		assert.InDelta(t, -148.0, data.(float64), 0.0)
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		require.Fail(t, " Temperature Sensor Event \"Data\" was not published")
	}

	assert.InDelta(t, -148.0, d.Value(), 0.0)
}
