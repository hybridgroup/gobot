package gpio

import (
	"errors"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

func TestAnalogSensorDriver(t *testing.T) {
	d := NewAnalogSensorDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
	gobot.Assert(t, d.Name(), "bot")
	gobot.Assert(t, d.Connection().Name(), "adaptor")

	d = NewAnalogSensorDriver(newGpioTestAdaptor("adaptor"), "bot", "1", 30*time.Second)
	gobot.Assert(t, d.interval, 30*time.Second)

	testAdaptorAnalogRead = func() (val int, err error) {
		val = 100
		return
	}
	ret := d.Command("Read")(nil).(map[string]interface{})

	gobot.Assert(t, ret["val"].(int), 100)
	gobot.Assert(t, ret["err"], nil)
}

func TestAnalogSensorDriverStart(t *testing.T) {
	sem := make(chan bool, 1)

	d := NewAnalogSensorDriver(newGpioTestAdaptor("adaptor"), "bot", "1")

	testAdaptorAnalogRead = func() (val int, err error) {
		val = 0
		return
	}
	gobot.Assert(t, len(d.Start()), 0)

	// data was received
	gobot.Once(d.Event(Data), func(data interface{}) {
		gobot.Assert(t, data.(int), 100)
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
	gobot.Once(d.Event(Error), func(data interface{}) {
		gobot.Assert(t, data.(error).Error(), "read error")
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
	gobot.Once(d.Event(Data), func(data interface{}) {
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
	gobot.Assert(t, len(d.Halt()), 0)
}
