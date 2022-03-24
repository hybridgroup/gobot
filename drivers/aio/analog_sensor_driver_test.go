package aio

import (
	"errors"
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*AnalogSensorDriver)(nil)

func TestAnalogSensorDriver(t *testing.T) {
	a := newAioTestAdaptor()
	d := NewAnalogSensorDriver(a, "1")
	gobottest.Refute(t, d.Connection(), nil)

	// default interval
	gobottest.Assert(t, d.interval, 10*time.Millisecond)

	// commands
	a = newAioTestAdaptor()
	d = NewAnalogSensorDriver(a, "42", 30*time.Second)
	d.SetScaler(func(input int) float64 { return 2.5*float64(input) - 3 })
	gobottest.Assert(t, d.Pin(), "42")
	gobottest.Assert(t, d.interval, 30*time.Second)

	a.TestAdaptorAnalogRead(func() (val int, err error) {
		val = 100
		return
	})
	ret := d.Command("Read")(nil).(map[string]interface{})
	gobottest.Assert(t, ret["val"].(int), 100)
	gobottest.Assert(t, ret["err"], nil)

	ret = d.Command("ReadValue")(nil).(map[string]interface{})
	gobottest.Assert(t, ret["val"].(float64), 247.0)
	gobottest.Assert(t, ret["err"], nil)

	// refresh value on read
	a = newAioTestAdaptor()
	d = NewAnalogSensorDriver(a, "3")
	a.TestAdaptorAnalogRead(func() (val int, err error) {
		val = 150
		return
	})
	gobottest.Assert(t, d.Value(), 0.0)
	val, err := d.ReadValue()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, 150.0)
	gobottest.Assert(t, d.Value(), 150.0)
	gobottest.Assert(t, d.RawValue(), 150)
}

func TestAnalogSensorDriverWithLinearScaler(t *testing.T) {
	// the input scales per default from 0...255
	var tests = map[string]struct {
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
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d.SetScaler(AnalogSensorLinearScaler(0, 255, tt.toMin, tt.toMax))
			a.TestAdaptorAnalogRead(func() (val int, err error) {
				return tt.input, nil
			})
			// act
			got, err := d.ReadValue()
			// assert
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, got, tt.want)
		})
	}
}

func TestAnalogSensorDriverStart(t *testing.T) {
	sem := make(chan bool, 1)
	a := newAioTestAdaptor()
	d := NewAnalogSensorDriver(a, "1")
	d.SetScaler(func(input int) float64 { return float64(input * input) })

	// expect data to be received
	d.Once(d.Event(Data), func(data interface{}) {
		gobottest.Assert(t, data.(int), 100)
		sem <- true
	})

	d.Once(d.Event(Value), func(data interface{}) {
		gobottest.Assert(t, data.(float64), 10000.0)
		sem <- true
	})

	// send data
	a.TestAdaptorAnalogRead(func() (val int, err error) {
		val = 100
		return
	})

	gobottest.Assert(t, d.Start(), nil)

	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf("AnalogSensor Event \"Data\" was not published")
	}

	// expect error to be received
	d.Once(d.Event(Error), func(data interface{}) {
		gobottest.Assert(t, data.(error).Error(), "read error")
		sem <- true
	})

	// send error
	a.TestAdaptorAnalogRead(func() (val int, err error) {
		err = errors.New("read error")
		return
	})

	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf("AnalogSensor Event \"Error\" was not published")
	}

	// send a halt message
	d.Once(d.Event(Data), func(data interface{}) {
		sem <- true
	})

	d.Once(d.Event(Value), func(data interface{}) {
		sem <- true
	})

	a.TestAdaptorAnalogRead(func() (val int, err error) {
		val = 200
		return
	})

	d.halt <- true

	select {
	case <-sem:
		t.Errorf("AnalogSensor Event should not published")
	case <-time.After(1 * time.Second):
	}
}

func TestAnalogSensorDriverHalt(t *testing.T) {
	d := NewAnalogSensorDriver(newAioTestAdaptor(), "1")
	done := make(chan struct{})
	go func() {
		<-d.halt
		close(done)
	}()
	gobottest.Assert(t, d.Halt(), nil)
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("AnalogSensor was not halted")
	}
}

func TestAnalogSensorDriverDefaultName(t *testing.T) {
	d := NewAnalogSensorDriver(newAioTestAdaptor(), "1")
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "AnalogSensor"), true)
}

func TestAnalogSensorDriverSetName(t *testing.T) {
	d := NewAnalogSensorDriver(newAioTestAdaptor(), "1")
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}
