package aio

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*GroveTemperatureSensorDriver)(nil)

func TestGroveTemperatureSensorDriver(t *testing.T) {
	testAdaptor := newAioTestAdaptor()
	d := NewGroveTemperatureSensorDriver(testAdaptor, "123")
	gobottest.Assert(t, d.Connection(), testAdaptor)
	gobottest.Assert(t, d.Pin(), "123")
	gobottest.Assert(t, d.interval, 10*time.Millisecond)
}

func TestGroveTemperatureSensorDriverScaling(t *testing.T) {
	var tests = map[string]struct {
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
			a.TestAdaptorAnalogRead(func() (val int, err error) {
				val = tt.input
				return
			})
			// act
			got, err := d.ReadValue()
			// assert
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, got, tt.want)
		})
	}
}

func TestGroveTempSensorPublishesTemperatureInCelsius(t *testing.T) {
	sem := make(chan bool, 1)
	a := newAioTestAdaptor()
	d := NewGroveTemperatureSensorDriver(a, "1")

	a.TestAdaptorAnalogRead(func() (val int, err error) {
		val = 585
		return
	})
	d.Once(d.Event(Value), func(data interface{}) {
		gobottest.Assert(t, fmt.Sprintf("%.2f", data.(float64)), "31.62")
		sem <- true
	})
	gobottest.Assert(t, d.Start(), nil)

	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf("Grove Temperature Sensor Event \"Data\" was not published")
	}

	gobottest.Assert(t, d.Temperature(), 31.61532462352477)
}

func TestGroveTempDriverDefaultName(t *testing.T) {
	d := NewGroveTemperatureSensorDriver(newAioTestAdaptor(), "1")
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "GroveTemperatureSensor"), true)
}
