//nolint:forcetypeassert // ok here
package aio

import (
	"errors"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

type groveDriverTestDriverAndEventer interface {
	gobot.Driver
	gobot.Eventer
}

func TestNewGroveRotaryDriver(t *testing.T) {
	// arrange
	a := newAioTestAdaptor()
	pin := "456"
	// act
	d := NewGroveRotaryDriver(a, pin)
	// assert: driver attributes
	assert.IsType(t, &GroveRotaryDriver{}, d)
	assert.NotNil(t, d.driverCfg)
	assert.True(t, strings.HasPrefix(d.Name(), "GroveRotary"))
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

func TestNewGroveLightSensorDriver(t *testing.T) {
	// arrange
	a := newAioTestAdaptor()
	pin := "456"
	// act
	d := NewGroveLightSensorDriver(a, pin)
	// assert: driver attributes
	assert.IsType(t, &GroveLightSensorDriver{}, d)
	assert.NotNil(t, d.driverCfg)
	assert.True(t, strings.HasPrefix(d.Name(), "GroveLightSensor"))
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

func TestNewGrovePiezoVibrationSensorDriver(t *testing.T) {
	// arrange
	a := newAioTestAdaptor()
	pin := "456"
	// act
	d := NewGrovePiezoVibrationSensorDriver(a, pin)
	// assert: driver attributes
	assert.IsType(t, &GrovePiezoVibrationSensorDriver{}, d)
	assert.NotNil(t, d.driverCfg)
	assert.True(t, strings.HasPrefix(d.Name(), "GrovePiezoVibrationSensor"))
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

func TestNewGroveSoundSensorDriver(t *testing.T) {
	// arrange
	a := newAioTestAdaptor()
	pin := "456"
	// act
	d := NewGroveSoundSensorDriver(a, pin)
	// assert: driver attributes
	assert.IsType(t, &GroveSoundSensorDriver{}, d)
	assert.NotNil(t, d.driverCfg)
	assert.True(t, strings.HasPrefix(d.Name(), "GroveSoundSensor"))
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

func TestGroveDriverHalt_WithSensorCyclicRead(t *testing.T) {
	// arrange
	testAdaptor := newAioTestAdaptor()
	pin := "456"

	drivers := []groveDriverTestDriverAndEventer{
		NewGroveSoundSensorDriver(testAdaptor, pin, WithSensorCyclicRead(10*time.Millisecond)),
		NewGroveLightSensorDriver(testAdaptor, pin, WithSensorCyclicRead(10*time.Millisecond)),
		NewGrovePiezoVibrationSensorDriver(testAdaptor, pin, WithSensorCyclicRead(10*time.Millisecond)),
		NewGroveRotaryDriver(testAdaptor, pin, WithSensorCyclicRead(10*time.Millisecond)),
	}

	for _, driver := range drivers {
		var callCount int32
		testAdaptor.analogReadFunc = func() (int, error) {
			atomic.AddInt32(&callCount, 1)
			return 42, nil
		}

		// Start the driver and allow for multiple digital reads
		_ = driver.Start()
		time.Sleep(20 * time.Millisecond)

		_ = driver.Halt()
		lastCallCount := atomic.LoadInt32(&callCount)
		// If driver was not halted, digital reads would still continue
		time.Sleep(20 * time.Millisecond)
		// note: if a reading is already in progress, it will be finished before halt have an impact
		if atomic.LoadInt32(&callCount) > lastCallCount+1 {
			require.Fail(t, "AnalogRead was called more than once after driver was halted")
		}
	}
}

func TestGroveDriverWithSensorCyclicReadPublishesError(t *testing.T) {
	// arrange
	testAdaptor := newAioTestAdaptor()
	pin := "456"

	drivers := []groveDriverTestDriverAndEventer{
		NewGroveSoundSensorDriver(testAdaptor, pin, WithSensorCyclicRead(10*time.Millisecond)),
		NewGroveLightSensorDriver(testAdaptor, pin, WithSensorCyclicRead(10*time.Millisecond)),
		NewGrovePiezoVibrationSensorDriver(testAdaptor, pin, WithSensorCyclicRead(10*time.Millisecond)),
		NewGroveRotaryDriver(testAdaptor, pin, WithSensorCyclicRead(10*time.Millisecond)),
	}

	for _, driver := range drivers {
		sem := make(chan struct{}, 1)
		// send error
		testAdaptor.analogReadFunc = func() (int, error) {
			return 0, errors.New("read error")
		}

		require.NoError(t, driver.Start())

		// expect error
		_ = driver.Once(driver.Event(Error), func(data interface{}) {
			assert.Equal(t, "read error", data.(error).Error())
			close(sem)
		})

		select {
		case <-sem:
		case <-time.After(time.Second):
			require.Fail(t, "%s Event \"Error\" was not published", groveGetType(driver))
		}

		// Cleanup
		_ = driver.Halt()
	}
}

func groveGetType(driver interface{}) string {
	d := reflect.TypeOf(driver)

	if d.Kind() == reflect.Ptr {
		return d.Elem().Name()
	}

	return d.Name()
}
