package aio

import (
	"errors"
	"fmt"
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

func TestGroveTempSensorPublishesTemperatureInCelsius(t *testing.T) {
	sem := make(chan bool, 1)
	d := NewGroveTemperatureSensorDriver(newAioTestAdaptor(), "1")

	testAdaptorAnalogRead = func() (val int, err error) {
		val = 585
		return
	}
	gobottest.Assert(t, d.Start(), nil)

	d.Once(d.Event(Data), func(data interface{}) {
		gobottest.Assert(t, fmt.Sprintf("%.2f", data.(float64)), "31.62")
		sem <- true
	})
}

func TestGroveTempSensorPublishesError(t *testing.T) {
	sem := make(chan bool, 1)
	d := NewGroveTemperatureSensorDriver(newAioTestAdaptor(), "1")

	// send error
	testAdaptorAnalogRead = func() (val int, err error) {
		err = errors.New("read error")
		return
	}

	gobottest.Assert(t, d.Start(), nil)

	// expect error
	d.Once(d.Event(Error), func(data interface{}) {
		gobottest.Assert(t, data.(error).Error(), "read error")
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(time.Second):
		t.Errorf("Grove Temperature Sensor Event \"Error\" was not published")
	}
}

func TestGroveTempSensorHalt(t *testing.T) {
	d := NewGroveTemperatureSensorDriver(newAioTestAdaptor(), "1")
	done := make(chan struct{})
	go func() {
		<-d.halt
		close(done)
	}()
	gobottest.Assert(t, d.Halt(), nil)
	select {
	case <-done:
	case <-time.After(time.Millisecond):
		t.Errorf("Grove Temperature Sensor was not halted")
	}
}
