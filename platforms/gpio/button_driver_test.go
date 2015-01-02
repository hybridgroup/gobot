package gpio

import (
	"errors"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

func initTestButtonDriver() *ButtonDriver {
	return NewButtonDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestButtonDriverHalt(t *testing.T) {
	d := initTestButtonDriver()
	go func() {
		<-d.halt
	}()
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestButtonDriver(t *testing.T) {
	d := NewButtonDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
	gobot.Assert(t, d.Name(), "bot")
	gobot.Assert(t, d.Connection().Name(), "adaptor")

	d = NewButtonDriver(newGpioTestAdaptor("adaptor"), "bot", "1", 30*time.Second)
	gobot.Assert(t, d.interval, 30*time.Second)
}

func TestButtonDriverStart(t *testing.T) {
	sem := make(chan bool, 0)
	d := initTestButtonDriver()
	gobot.Assert(t, len(d.Start()), 0)

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}

	gobot.Once(d.Event(Push), func(data interface{}) {
		gobot.Assert(t, d.Active, true)
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(15 * time.Millisecond):
		t.Errorf("Button Event \"Push\" was not published")
	}

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 0
		return
	}

	gobot.Once(d.Event(Release), func(data interface{}) {
		gobot.Assert(t, d.Active, false)
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(15 * time.Millisecond):
		t.Errorf("Button Event \"Release\" was not published")
	}

	testAdaptorDigitalRead = func() (val int, err error) {
		err = errors.New("digital read error")
		return
	}

	gobot.Once(d.Event(Error), func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(15 * time.Millisecond):
		t.Errorf("Button Event \"Error\" was not published")
	}

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}

	gobot.Once(d.Event(Push), func(data interface{}) {
		sem <- true
	})

	d.halt <- true

	select {
	case <-sem:
		t.Errorf("Button Event \"Press\" should not published")
	case <-time.After(15 * time.Millisecond):
	}

}
