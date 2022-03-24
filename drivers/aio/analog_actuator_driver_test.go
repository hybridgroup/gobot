package aio

import (
	"strings"
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestAnalogActuatorDriver(t *testing.T) {
	a := newAioTestAdaptor()
	d := NewAnalogActuatorDriver(a, "47")

	gobottest.Refute(t, d.Connection(), nil)
	gobottest.Assert(t, d.Pin(), "47")

	err := d.RawWrite(100)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.written), 1)
	gobottest.Assert(t, a.written[0], 100)

	err = d.Write(247.0)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.written), 2)
	gobottest.Assert(t, a.written[1], 247)
	gobottest.Assert(t, d.RawValue(), 247)
	gobottest.Assert(t, d.Value(), 247.0)
}

func TestAnalogActuatorDriverWithScaler(t *testing.T) {
	// commands
	a := newAioTestAdaptor()
	d := NewAnalogActuatorDriver(a, "7")
	d.SetScaler(func(input float64) int { return int((input + 3) / 2.5) })

	err := d.Command("RawWrite")(map[string]interface{}{"val": "100"})
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.written), 1)
	gobottest.Assert(t, a.written[0], 100)

	err = d.Command("Write")(map[string]interface{}{"val": "247.0"})
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.written), 2)
	gobottest.Assert(t, a.written[1], 100)
}

func TestAnalogActuatorDriverLinearScaler(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, len(a.written), 1)
			gobottest.Assert(t, a.written[0], tt.want)
		})
	}
}

func TestAnalogActuatorDriverStart(t *testing.T) {
	d := NewAnalogActuatorDriver(newAioTestAdaptor(), "1")
	gobottest.Assert(t, d.Start(), nil)
}

func TestAnalogActuatorDriverHalt(t *testing.T) {
	d := NewAnalogActuatorDriver(newAioTestAdaptor(), "1")
	gobottest.Assert(t, d.Halt(), nil)
}

func TestAnalogActuatorDriverDefaultName(t *testing.T) {
	d := NewAnalogActuatorDriver(newAioTestAdaptor(), "1")
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "AnalogActuator"), true)
}

func TestAnalogActuatorDriverSetName(t *testing.T) {
	d := NewAnalogActuatorDriver(newAioTestAdaptor(), "1")
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}
