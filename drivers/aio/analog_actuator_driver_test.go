//nolint:forcetypeassert // ok here
package aio

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAnalogActuatorDriver(t *testing.T) {
	// arrange
	const pin = "47"
	a := newAioTestAdaptor()
	// act
	d := NewAnalogActuatorDriver(a, pin)
	// assert: driver attributes
	assert.IsType(t, &AnalogActuatorDriver{}, d)
	assert.NotNil(t, d.driverCfg)
	assert.True(t, strings.HasPrefix(d.Name(), "AnalogActuator"))
	assert.Equal(t, a, d.Connection())
	require.NoError(t, d.afterStart())
	require.NoError(t, d.beforeHalt())
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: actuator attributes
	assert.Equal(t, pin, d.Pin())
	assert.InDelta(t, 0.0, d.lastValue, 0, 0)
	assert.Equal(t, 0, d.lastRawValue)
	require.NotNil(t, d.actuatorCfg)
	assert.NotNil(t, d.actuatorCfg.scale)
}

func TestNewAnalogActuatorDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	myName := "relay 1"
	myScaler := func(input float64) int { return int(2 * input) }
	panicFunc := func() {
		NewAnalogActuatorDriver(newAioTestAdaptor(), "1", WithName("crazy"), WithSensorCyclicRead(10*time.Millisecond))
	}
	// act
	d := NewAnalogActuatorDriver(newAioTestAdaptor(), "1", WithName(myName), WithActuatorScaler(myScaler))
	// assert
	assert.Equal(t, myName, d.Name())
	assert.Equal(t, 3, d.actuatorCfg.scale(1.5))
	assert.PanicsWithValue(t, "'read interval option for analog sensors' can not be applied on 'crazy'", panicFunc)
}

func TestAnalogActuatorWriteRaw(t *testing.T) {
	tests := map[string]struct {
		inputVal         int
		simulateWriteErr bool
		wantWritten      int
		wantErr          string
	}{
		"write_raw":   {inputVal: 100, wantWritten: 100},
		"error_write": {inputVal: 12345, wantWritten: 12345, simulateWriteErr: true, wantErr: "write error"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			const pin = "47"
			a := newAioTestAdaptor()
			d := NewAnalogActuatorDriver(a, pin)
			a.simulateWriteError = tc.simulateWriteErr
			a.written = nil // reset previous writes
			// act
			err := d.WriteRaw(tc.inputVal)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				assert.Empty(t, a.written)
			} else {
				require.NoError(t, err)
				assert.Len(t, a.written, 1)
				assert.Equal(t, pin, a.written[0].pin)
				assert.Equal(t, tc.wantWritten, a.written[0].val)
			}
		})
	}
}

func TestAnalogActuatorWriteRaw_AnalogWriteNotSupported(t *testing.T) {
	// arrange
	d := NewAnalogActuatorDriver(newAioTestAdaptor(), "1")
	d.connection = &aioTestBareAdaptor{}
	// act & assert
	require.EqualError(t, d.WriteRaw(3), "AnalogWrite is not supported by the platform 'bare'")
}

func TestAnalogActuatorWrite_SetScaler(t *testing.T) {
	tests := map[string]struct {
		fromMin     float64
		fromMax     float64
		input       float64
		wantWritten int
	}{
		"byte_range_min":           {fromMin: 0, fromMax: 255, input: 0, wantWritten: 0},
		"byte_range_max":           {fromMin: 0, fromMax: 255, input: 255, wantWritten: 255},
		"signed_percent_range_min": {fromMin: -100, fromMax: 100, input: -100, wantWritten: 0},
		"signed_percent_range_mid": {fromMin: -100, fromMax: 100, input: 0, wantWritten: 127},
		"signed_percent_range_max": {fromMin: -100, fromMax: 100, input: 100, wantWritten: 255},
		"voltage_range_min":        {fromMin: 0, fromMax: 5.1, input: 0, wantWritten: 0},
		"voltage_range_nearmin":    {fromMin: 0, fromMax: 5.1, input: 0.02, wantWritten: 1},
		"voltage_range_mid":        {fromMin: 0, fromMax: 5.1, input: 2.55, wantWritten: 127},
		"voltage_range_nearmax":    {fromMin: 0, fromMax: 5.1, input: 5.08, wantWritten: 254},
		"voltage_range_max":        {fromMin: 0, fromMax: 5.1, input: 5.1, wantWritten: 255},
		"upscale":                  {fromMin: 0, fromMax: 24, input: 12, wantWritten: 127},
		"below_min":                {fromMin: -10, fromMax: 10, input: -11, wantWritten: 0},
		"exceed_max":               {fromMin: 0, fromMax: 20, input: 21, wantWritten: 255},
	}

	const pin = "7"
	a := newAioTestAdaptor()
	d := NewAnalogActuatorDriver(a, pin)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d.SetScaler(AnalogActuatorLinearScaler(tc.fromMin, tc.fromMax, 0, 255))
			a.written = nil // reset previous writes
			// act
			err := d.Write(tc.input)
			// assert
			require.NoError(t, err)
			assert.Len(t, a.written, 1)
			assert.Equal(t, pin, a.written[0].pin)
			assert.Equal(t, tc.wantWritten, a.written[0].val)
		})
	}
}

func TestAnalogActuatorCommands_WithActuatorScaler(t *testing.T) {
	// arrange
	const pin = "8"
	a := newAioTestAdaptor()
	d := NewAnalogActuatorDriver(a, pin, WithActuatorScaler(func(input float64) int { return int((input + 3) / 2.5) }))
	a.written = nil // reset previous writes
	// act & assert: WriteRaw
	err := d.Command("WriteRaw")(map[string]interface{}{"val": "100"})
	assert.Nil(t, err)
	assert.Len(t, a.written, 1)
	assert.Equal(t, pin, a.written[0].pin)
	assert.Equal(t, 100, a.written[0].val)
	assert.Equal(t, 100, d.RawValue())
	assert.InDelta(t, 0.0, d.Value(), 0.0)
	// act & assert: Write
	err = d.Command("Write")(map[string]interface{}{"val": "247.0"})
	assert.Nil(t, err)
	assert.Len(t, a.written, 2)
	assert.Equal(t, pin, a.written[1].pin)
	assert.Equal(t, 100, a.written[1].val)
	assert.Equal(t, 100, d.RawValue())
	assert.InDelta(t, 247.0, d.Value(), 0.0)
	// arrange & act & assert: Write with error
	a.simulateWriteError = true
	err = d.Command("Write")(map[string]interface{}{"val": "247.0"})
	require.EqualError(t, err.(error), "write error")
}
