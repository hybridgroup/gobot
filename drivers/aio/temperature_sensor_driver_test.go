package aio

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTemperatureSensorDriver(t *testing.T) {
	testAdaptor := newAioTestAdaptor()
	d := NewTemperatureSensorDriver(testAdaptor, "123")
	assert.Equal(t, testAdaptor, d.Connection())
	assert.Equal(t, "123", d.Pin())
	assert.Equal(t, 10*time.Millisecond, d.interval)
}

func TestTemperatureSensorDriverNtcScaling(t *testing.T) {
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
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a.analogReadFunc = func() (val int, err error) {
				val = tt.input
				return
			}
			// act
			got, err := d.Read()
			// assert
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTemperatureSensorDriverLinearScaling(t *testing.T) {
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
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a.analogReadFunc = func() (val int, err error) {
				val = tt.input
				return
			}
			// act
			got, err := d.Read()
			// assert
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTempSensorPublishesTemperatureInCelsius(t *testing.T) {
	sem := make(chan bool, 1)
	a := newAioTestAdaptor()
	d := NewTemperatureSensorDriver(a, "1")
	ntc := TemperatureSensorNtcConf{TC0: 25, R0: 10000.0, B: 3975} // Ohm, R25=10k
	d.SetNtcScaler(1023, 10000, false, ntc)                        // Ohm, reference value: 1023, series R: 10k

	a.analogReadFunc = func() (val int, err error) {
		val = 585
		return
	}
	_ = d.Once(d.Event(Value), func(data interface{}) {
		assert.Equal(t, "31.62", fmt.Sprintf("%.2f", data.(float64)))
		sem <- true
	})
	assert.NoError(t, d.Start())

	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf(" Temperature Sensor Event \"Data\" was not published")
	}

	assert.Equal(t, 31.61532462352477, d.Value())
}

func TestTempSensorPublishesError(t *testing.T) {
	sem := make(chan bool, 1)
	a := newAioTestAdaptor()
	d := NewTemperatureSensorDriver(a, "1")

	// send error
	a.analogReadFunc = func() (val int, err error) {
		err = errors.New("read error")
		return
	}

	assert.NoError(t, d.Start())

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

func TestTempSensorHalt(t *testing.T) {
	d := NewTemperatureSensorDriver(newAioTestAdaptor(), "1")
	done := make(chan struct{})
	go func() {
		<-d.halt
		close(done)
	}()
	assert.NoError(t, d.Halt())
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Errorf(" Temperature Sensor was not halted")
	}
}

func TestTempDriverDefaultName(t *testing.T) {
	d := NewTemperatureSensorDriver(newAioTestAdaptor(), "1")
	assert.True(t, strings.HasPrefix(d.Name(), "TemperatureSensor"))
}

func TestTempDriverSetName(t *testing.T) {
	d := NewTemperatureSensorDriver(newAioTestAdaptor(), "1")
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
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
			want:  TemperatureSensorNtcConf{TC0: 25, R0: 2500.0, TC1: -13, R1: 10000, B: 2829.6355560320544, t0: 298.15, r: 9.490644159087891},
		},
		"T1_high": {
			input: TemperatureSensorNtcConf{TC0: 25, R0: 2500.0, TC1: 100, R1: 371},
			want:  TemperatureSensorNtcConf{TC0: 25, R0: 2500.0, TC1: 100, R1: 371, B: 2830.087381913779, t0: 298.15, r: 9.49215959052081},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			ntc := tt.input
			// act
			ntc.initialize()
			// assert
			assert.Equal(t, tt.want, ntc)
		})
	}
}
