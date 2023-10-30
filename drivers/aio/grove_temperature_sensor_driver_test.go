package aio

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*GroveTemperatureSensorDriver)(nil)

func TestGroveTemperatureSensorDriver(t *testing.T) {
	testAdaptor := newAioTestAdaptor()
	d := NewGroveTemperatureSensorDriver(testAdaptor, "123")
	assert.Equal(t, testAdaptor, d.Connection())
	assert.Equal(t, "123", d.Pin())
	assert.Equal(t, 10*time.Millisecond, d.interval)
}

func TestGroveTemperatureSensorDriverScaling(t *testing.T) {
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

func TestGroveTempSensorPublishesTemperatureInCelsius(t *testing.T) {
	sem := make(chan bool, 1)
	a := newAioTestAdaptor()
	d := NewGroveTemperatureSensorDriver(a, "1")

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
		t.Errorf("Grove Temperature Sensor Event \"Data\" was not published")
	}

	assert.Equal(t, 31.61532462352477, d.Temperature())
}

func TestGroveTempDriverDefaultName(t *testing.T) {
	d := NewGroveTemperatureSensorDriver(newAioTestAdaptor(), "1")
	assert.True(t, strings.HasPrefix(d.Name(), "GroveTemperatureSensor"))
}
