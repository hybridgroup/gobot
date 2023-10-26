package aio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalogActuatorDriver(t *testing.T) {
	a := newAioTestAdaptor()
	d := NewAnalogActuatorDriver(a, "47")

	assert.NotNil(t, d.Connection())
	assert.Equal(t, "47", d.Pin())

	err := d.RawWrite(100)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(a.written))
	assert.Equal(t, 100, a.written[0])

	err = d.Write(247.0)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(a.written))
	assert.Equal(t, 247, a.written[1])
	assert.Equal(t, 247, d.RawValue())
	assert.Equal(t, 247.0, d.Value())
}

func TestAnalogActuatorDriverWithScaler(t *testing.T) {
	// commands
	a := newAioTestAdaptor()
	d := NewAnalogActuatorDriver(a, "7")
	d.SetScaler(func(input float64) int { return int((input + 3) / 2.5) })

	err := d.Command("RawWrite")(map[string]interface{}{"val": "100"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(a.written))
	assert.Equal(t, 100, a.written[0])

	err = d.Command("Write")(map[string]interface{}{"val": "247.0"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(a.written))
	assert.Equal(t, 100, a.written[1])
}

func TestAnalogActuatorDriverLinearScaler(t *testing.T) {
	tests := map[string]struct {
		fromMin float64
		fromMax float64
		input   float64
		want    int
	}{
		"byte_range_min":           {fromMin: 0, fromMax: 255, input: 0, want: 0},
		"byte_range_max":           {fromMin: 0, fromMax: 255, input: 255, want: 255},
		"signed_percent_range_min": {fromMin: -100, fromMax: 100, input: -100, want: 0},
		"signed_percent_range_mid": {fromMin: -100, fromMax: 100, input: 0, want: 127},
		"signed_percent_range_max": {fromMin: -100, fromMax: 100, input: 100, want: 255},
		"voltage_range_min":        {fromMin: 0, fromMax: 5.1, input: 0, want: 0},
		"voltage_range_nearmin":    {fromMin: 0, fromMax: 5.1, input: 0.02, want: 1},
		"voltage_range_mid":        {fromMin: 0, fromMax: 5.1, input: 2.55, want: 127},
		"voltage_range_nearmax":    {fromMin: 0, fromMax: 5.1, input: 5.08, want: 254},
		"voltage_range_max":        {fromMin: 0, fromMax: 5.1, input: 5.1, want: 255},
		"upscale":                  {fromMin: 0, fromMax: 24, input: 12, want: 127},
		"below_min":                {fromMin: -10, fromMax: 10, input: -11, want: 0},
		"exceed_max":               {fromMin: 0, fromMax: 20, input: 21, want: 255},
	}
	a := newAioTestAdaptor()
	d := NewAnalogActuatorDriver(a, "7")

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d.SetScaler(AnalogActuatorLinearScaler(tt.fromMin, tt.fromMax, 0, 255))
			a.written = []int{} // reset previous write
			// act
			err := d.Write(tt.input)
			// assert
			assert.NoError(t, err)
			assert.Equal(t, 1, len(a.written))
			assert.Equal(t, tt.want, a.written[0])
		})
	}
}

func TestAnalogActuatorDriverStart(t *testing.T) {
	d := NewAnalogActuatorDriver(newAioTestAdaptor(), "1")
	assert.NoError(t, d.Start())
}

func TestAnalogActuatorDriverHalt(t *testing.T) {
	d := NewAnalogActuatorDriver(newAioTestAdaptor(), "1")
	assert.NoError(t, d.Halt())
}

func TestAnalogActuatorDriverDefaultName(t *testing.T) {
	d := NewAnalogActuatorDriver(newAioTestAdaptor(), "1")
	assert.True(t, strings.HasPrefix(d.Name(), "AnalogActuator"))
}

func TestAnalogActuatorDriverSetName(t *testing.T) {
	d := NewAnalogActuatorDriver(newAioTestAdaptor(), "1")
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}
