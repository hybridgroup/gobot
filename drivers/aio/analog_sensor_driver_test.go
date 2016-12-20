package aio

import (
	"errors"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*AnalogSensorDriver)(nil)

func TestAnalogSensorDriver(t *testing.T) {
	d := NewAnalogSensorDriver(newAioTestAdaptor(), "1")
	gobottest.Refute(t, d.Connection(), nil)
	// default interval
	gobottest.Assert(t, d.interval, 10*time.Millisecond)

	d = NewAnalogSensorDriver(newAioTestAdaptor(), "42", 30*time.Second)
	gobottest.Assert(t, d.Pin(), "42")
	gobottest.Assert(t, d.interval, 30*time.Second)

	testAdaptorAnalogRead = func() (val int, err error) {
		val = 100
		return
	}
	ret := d.Command("Read")(nil).(map[string]interface{})

	gobottest.Assert(t, ret["val"].(int), 100)
	gobottest.Assert(t, ret["err"], nil)
}

func TestAnalogSensorDriverStart(t *testing.T) {
	sem := make(chan bool, 1)

	d := NewAnalogSensorDriver(newAioTestAdaptor(), "1")

	gobottest.Assert(t, d.Start(), nil)

	// expect data to be received
	d.Once(d.Event(Data), func(data interface{}) {
		gobottest.Assert(t, data.(int), 100)
		sem <- true
	})

	// send data
	testAdaptorAnalogRead = func() (val int, err error) {
		val = 100
		return
	}

	select {
	case <-sem:
	case <-time.After(10 * time.Second):
		t.Errorf("AnalogSensor Event \"Data\" was not published")
	}

	// expect error to be received
	d.Once(d.Event(Error), func(data interface{}) {
		gobottest.Assert(t, data.(error).Error(), "read error")
		sem <- true
	})

	// send error
	testAdaptorAnalogRead = func() (val int, err error) {
		err = errors.New("read error")
		return
	}

	select {
	case <-sem:
	case <-time.After(10 * time.Second):
		t.Errorf("AnalogSensor Event \"Error\" was not published")
	}

	// send a halt message
	d.Once(d.Event(Data), func(data interface{}) {
		sem <- true
	})

	testAdaptorAnalogRead = func() (val int, err error) {
		val = 200
		return
	}

	d.halt <- true

	select {
	case <-sem:
		t.Errorf("AnalogSensor Event should not published")
	case <-time.After(100 * time.Millisecond):
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
	case <-time.After(time.Millisecond):
		t.Errorf("AnalogSensor was not halted")
	}
}
