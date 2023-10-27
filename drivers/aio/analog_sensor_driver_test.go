package aio

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*AnalogSensorDriver)(nil)

func TestAnalogSensorDriver(t *testing.T) {
	a := newAioTestAdaptor()
	d := NewAnalogSensorDriver(a, "1")
	assert.NotNil(t, d.Connection())

	// default interval
	assert.Equal(t, 10*time.Millisecond, d.interval)

	// commands
	a = newAioTestAdaptor()
	d = NewAnalogSensorDriver(a, "42", 30*time.Second)
	d.SetScaler(func(input int) float64 { return 2.5*float64(input) - 3 })
	assert.Equal(t, "42", d.Pin())
	assert.Equal(t, 30*time.Second, d.interval)

	a.analogReadFunc = func() (val int, err error) {
		val = 100
		return
	}

	ret := d.Command("ReadRaw")(nil).(map[string]interface{})
	assert.Equal(t, 100, ret["val"].(int))
	assert.Nil(t, ret["err"])

	ret = d.Command("Read")(nil).(map[string]interface{})
	assert.Equal(t, 247.0, ret["val"].(float64))
	assert.Nil(t, ret["err"])

	// refresh value on read
	a = newAioTestAdaptor()
	d = NewAnalogSensorDriver(a, "3")
	a.analogReadFunc = func() (val int, err error) {
		val = 150
		return
	}
	assert.Equal(t, 0.0, d.Value())
	val, err := d.Read()
	assert.NoError(t, err)
	assert.Equal(t, 150.0, val)
	assert.Equal(t, 150.0, d.Value())
	assert.Equal(t, 150, d.RawValue())
}

func TestAnalogSensorDriverWithLinearScaler(t *testing.T) {
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
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d.SetScaler(AnalogSensorLinearScaler(0, 255, tt.toMin, tt.toMax))
			a.analogReadFunc = func() (val int, err error) {
				return tt.input, nil
			}
			// act
			got, err := d.Read()
			// assert
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAnalogSensorDriverStart(t *testing.T) {
	sem := make(chan bool, 1)
	a := newAioTestAdaptor()
	d := NewAnalogSensorDriver(a, "1")
	d.SetScaler(func(input int) float64 { return float64(input * input) })

	// expect data to be received
	_ = d.Once(d.Event(Data), func(data interface{}) {
		assert.Equal(t, 100, data.(int))
		sem <- true
	})

	_ = d.Once(d.Event(Value), func(data interface{}) {
		assert.Equal(t, 10000.0, data.(float64))
		sem <- true
	})

	// send data
	a.analogReadFunc = func() (val int, err error) {
		val = 100
		return
	}

	assert.NoError(t, d.Start())

	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf("AnalogSensor Event \"Data\" was not published")
	}

	// expect error to be received
	_ = d.Once(d.Event(Error), func(data interface{}) {
		assert.Equal(t, "read error", data.(error).Error())
		sem <- true
	})

	// send error
	a.analogReadFunc = func() (val int, err error) {
		err = errors.New("read error")
		return
	}

	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf("AnalogSensor Event \"Error\" was not published")
	}

	// send a halt message
	_ = d.Once(d.Event(Data), func(data interface{}) {
		sem <- true
	})

	_ = d.Once(d.Event(Value), func(data interface{}) {
		sem <- true
	})

	a.analogReadFunc = func() (val int, err error) {
		val = 200
		return
	}

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
	assert.NoError(t, d.Halt())
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("AnalogSensor was not halted")
	}
}

func TestAnalogSensorDriverDefaultName(t *testing.T) {
	d := NewAnalogSensorDriver(newAioTestAdaptor(), "1")
	assert.True(t, strings.HasPrefix(d.Name(), "AnalogSensor"))
}

func TestAnalogSensorDriverSetName(t *testing.T) {
	d := NewAnalogSensorDriver(newAioTestAdaptor(), "1")
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}
