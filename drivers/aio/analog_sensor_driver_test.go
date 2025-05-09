//nolint:forcetypeassert // ok here
package aio

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*AnalogSensorDriver)(nil)

func TestNewAnalogSensorDriver(t *testing.T) {
	// arrange
	const pin = "5"
	a := newAioTestAdaptor()
	// act
	d := NewAnalogSensorDriver(a, pin)
	// assert: driver attributes
	assert.IsType(t, &AnalogSensorDriver{}, d)
	assert.NotNil(t, d.driverCfg)
	assert.True(t, strings.HasPrefix(d.Name(), "AnalogSensor"))
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

func TestNewAnalogSensorDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName     = "voltage 1"
		cycReadDur = 10 * time.Millisecond
	)
	panicFunc := func() {
		NewAnalogSensorDriver(newAioTestAdaptor(), "1", WithName("crazy"), WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewAnalogSensorDriver(newAioTestAdaptor(), "1", WithName(myName), WithSensorCyclicRead(cycReadDur))
	// assert
	assert.Equal(t, cycReadDur, d.sensorCfg.readInterval)
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestAnalogSensor_WithSensorScaler(t *testing.T) {
	// arrange
	myScaler := func(input int) float64 { return float64(input) / 2 }
	cfg := sensorConfiguration{}
	// act
	WithSensorScaler(myScaler).apply(&cfg)
	// assert
	assert.InDelta(t, 1.5, cfg.scale(3), 0.0)
}

func TestAnalogSensorDriverReadRaw(t *testing.T) {
	tests := map[string]struct {
		simulateReadErr bool
		wantVal         int
		wantErr         string
	}{
		"read_raw":   {wantVal: analogReadReturnValue},
		"error_read": {wantVal: 0, simulateReadErr: true, wantErr: "read error"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			const pin = "47"
			a := newAioTestAdaptor()
			d := NewAnalogSensorDriver(a, pin)
			a.simulateReadError = tc.simulateReadErr
			a.written = nil // reset previous writes
			// act
			got, err := d.ReadRaw()
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				assert.Empty(t, a.written)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantVal, got)
		})
	}
}

func TestAnalogSensorDriverReadRaw_AnalogWriteNotSupported(t *testing.T) {
	// arrange
	d := NewAnalogSensorDriver(newAioTestAdaptor(), "1")
	d.connection = &aioTestBareAdaptor{}
	// act & assert
	got, err := d.ReadRaw()
	require.EqualError(t, err, "AnalogRead is not supported by the platform 'bare'")
	assert.Equal(t, 0, got)
}

func TestAnalogSensorRead_SetScaler(t *testing.T) {
	// the input scales per default from 0...255
	tests := map[string]struct {
		toMin float64
		toMax float64
		input int
		want  float64
	}{
		"single_byte_range_min":   {toMin: 0, toMax: 255, input: 0, want: 0},
		"single_byte_range_max":   {toMin: 0, toMax: 255, input: 255, want: 255},
		"single_below_min":        {toMin: 3, toMax: 121, input: -1, want: 3},
		"single_is_max":           {toMin: 5, toMax: 6, input: 255, want: 6},
		"single_upscale":          {toMin: 337, toMax: 5337, input: 127, want: 2827.196078431373},
		"grd_int_range_min":       {toMin: -180, toMax: 180, input: 0, want: -180},
		"grd_int_range_minus_one": {toMin: -180, toMax: 180, input: 127, want: -0.7058823529411598},
		"grd_int_range_max":       {toMin: -180, toMax: 180, input: 255, want: 180},
		"upscale":                 {toMin: -10, toMax: 1234, input: 255, want: 1234},
	}
	a := newAioTestAdaptor()
	d := NewAnalogSensorDriver(a, "7")
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d.SetScaler(AnalogSensorLinearScaler(0, 255, tc.toMin, tc.toMax))
			a.analogReadFunc = func() (int, error) {
				return tc.input, nil
			}
			// act
			got, err := d.Read()
			// assert
			require.NoError(t, err)
			assert.InDelta(t, tc.want, got, 1.0e-14)
		})
	}
}

func TestAnalogSensor_WithSensorCyclicRead(t *testing.T) {
	// arrange
	a := newAioTestAdaptor()
	d := NewAnalogSensorDriver(a, "1", WithSensorCyclicRead(10*time.Millisecond))
	d.SetScaler(func(input int) float64 { return float64(input * input) })
	semData := make(chan bool)
	semDone := make(chan bool)
	nextVal := make(chan int)
	readTimeout := time.Second
	a.analogReadFunc = func() (int, error) {
		val := 100
		var err error
		select {
		case val = <-nextVal:
			if val < 0 {
				err = fmt.Errorf("analog read error")
			}
			return val, err
		default:
			return val, nil
		}
	}

	// arrange: expect raw value to be received
	_ = d.Once(Data, func(data interface{}) { // we can't use d.Event(Data) here, because not registered yet
		assert.Equal(t, 100, data.(int))
		semData <- true
	})

	// arrange: expect scaled value to be received
	_ = d.Once(Value, func(value interface{}) { // we can't use d.Event(Value) here, because not registered yet
		assert.InDelta(t, 10000.0, value.(float64), 0.0)
		<-semData // wait for data is finished
		semDone <- true
		nextVal <- -1 // arrange: error in read function
	})

	// wait some time to ensure the cyclic go routine is working
	time.Sleep(15 * time.Millisecond)

	// act (start cyclic reading)
	require.NoError(t, d.Start())

	// assert: both events within timeout
	select {
	case <-semDone:
	case <-time.After(readTimeout):
		require.Fail(t, "AnalogSensor Event \"Data\" was not published")
	}

	// arrange: for error to be received
	_ = d.Once(d.Event(Error), func(err interface{}) {
		assert.Equal(t, "analog read error", err.(error).Error())
		semDone <- true
	})

	// assert: error
	select {
	case <-semDone:
	case <-time.After(readTimeout):
		require.Fail(t, "AnalogSensor Event \"Error\" was not published")
	}

	// arrange: for halt message
	_ = d.Once(d.Event(Data), func(data interface{}) {
		semData <- true
	})

	_ = d.Once(d.Event(Value), func(value interface{}) {
		semDone <- true
	})

	// act: send the halt message
	require.NoError(t, d.Halt())

	// assert: no event
	select {
	case <-semData:
		require.Fail(t, "AnalogSensor Event for data should not published")
	case <-semDone:
		require.Fail(t, "AnalogSensor Event for value should not published")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestAnalogSensorHalt_WithSensorCyclicRead(t *testing.T) {
	// arrange
	d := NewAnalogSensorDriver(newAioTestAdaptor(), "1", WithSensorCyclicRead(10*time.Millisecond))
	require.NoError(t, d.Start())
	timeout := 2 * d.sensorCfg.readInterval
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-d.halt: // wait until halt is broadcasted by close the channel
		case <-time.After(timeout): // otherwise run into the timeout
			assert.Fail(t, "halt was not received within %s", timeout)
		}
	}()
	// act & assert
	require.NoError(t, d.Halt())
	wg.Wait() // wait until the go function was really finished
}

func TestAnalogSensorCommands_WithSensorScaler(t *testing.T) {
	// arrange
	a := newAioTestAdaptor()
	d := NewAnalogSensorDriver(a, "42", WithSensorScaler(func(input int) float64 { return 2.5*float64(input) - 3 }))
	var readReturn int
	a.analogReadFunc = func() (int, error) {
		readReturn += 100
		return readReturn, nil
	}
	// act & assert: ReadRaw
	ret := d.Command("ReadRaw")(nil).(map[string]interface{})
	assert.Equal(t, 100, ret["val"].(int))
	assert.Nil(t, ret["err"])
	assert.Equal(t, 100, d.RawValue())
	assert.InDelta(t, 247.0, d.Value(), 0.0)
	// act & assert: Read
	ret = d.Command("Read")(nil).(map[string]interface{})
	assert.InDelta(t, 497.0, ret["val"].(float64), 0.0)
	assert.Nil(t, ret["err"])
	assert.Equal(t, 200, d.RawValue())
	assert.InDelta(t, 497.0, d.Value(), 0.0)
}
