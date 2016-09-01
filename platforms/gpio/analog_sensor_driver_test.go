package gpio

import (
	"errors"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Driver = (*AnalogSensorDriver)(nil)

func TestAnalogSensorDriver(t *testing.T) {
	d := NewAnalogSensorDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
	gobottest.Assert(t, d.Name(), "bot")
	gobottest.Assert(t, d.Connection().Name(), "adaptor")

	d = NewAnalogSensorDriver(newGpioTestAdaptor("adaptor"), "bot", "1", 30*time.Second)
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

	d := NewAnalogSensorDriver(newGpioTestAdaptor("adaptor"), "bot", "1")

	testAdaptorAnalogRead = func() (val int, err error) {
		val = 0
		return
	}
	gobottest.Assert(t, len(d.Start()), 0)

	// data was received
	d.Once(d.Event(Data), func(data interface{}) {
		gobottest.Assert(t, data.(int), 100)
		sem <- true
	})

	testAdaptorAnalogRead = func() (val int, err error) {
		val = 100
		return
	}

	select {
	case <-sem:
	case <-time.After(10 * time.Second):
		t.Errorf("AnalogSensor Event \"Data\" was not published")
	}

	// read error
	d.Once(d.Event(Error), func(data interface{}) {
		gobottest.Assert(t, data.(error).Error(), "read error")
		sem <- true
	})

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
	d := NewAnalogSensorDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
	go func() {
		<-d.halt
	}()
	gobottest.Assert(t, len(d.Halt()), 0)
}
