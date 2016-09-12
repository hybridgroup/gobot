package gpio

import (
	"errors"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Driver = (*ButtonDriver)(nil)

const BUTTON_TEST_DELAY = 150

func initTestButtonDriver() *ButtonDriver {
	return NewButtonDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestButtonDriverHalt(t *testing.T) {
	d := initTestButtonDriver()
	go func() {
		<-d.halt
	}()
	gobottest.Assert(t, len(d.Halt()), 0)
}

func TestButtonDriver(t *testing.T) {
	d := NewButtonDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
	gobottest.Assert(t, d.Name(), "bot")
	gobottest.Assert(t, d.Connection().Name(), "adaptor")

	d = NewButtonDriver(newGpioTestAdaptor("adaptor"), "bot", "1", 30*time.Second)
	gobottest.Assert(t, d.interval, 30*time.Second)
}

func TestButtonDriverStart(t *testing.T) {
	sem := make(chan bool, 0)
	d := initTestButtonDriver()
	gobottest.Assert(t, len(d.Start()), 0)

	d.Once(ButtonPush, func(data interface{}) {
		gobottest.Assert(t, d.Active, true)
		sem <- true
	})

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}

	select {
	case <-sem:
	case <-time.After(BUTTON_TEST_DELAY * time.Millisecond):
		t.Errorf("Button Event \"Push\" was not published")
	}

	d.Once(ButtonRelease, func(data interface{}) {
		gobottest.Assert(t, d.Active, false)
		sem <- true
	})

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 0
		return
	}

	select {
	case <-sem:
	case <-time.After(BUTTON_TEST_DELAY * time.Millisecond):
		t.Errorf("Button Event \"Release\" was not published")
	}

	testAdaptorDigitalRead = func() (val int, err error) {
		err = errors.New("digital read error")
		return
	}

	d.Once(Error, func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(BUTTON_TEST_DELAY * time.Millisecond):
		t.Errorf("Button Event \"Error\" was not published")
	}

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}

	d.Once(ButtonPush, func(data interface{}) {
		sem <- true
	})

	d.halt <- true

	select {
	case <-sem:
		t.Errorf("Button Event \"Press\" should not published")
	case <-time.After(BUTTON_TEST_DELAY * time.Millisecond):
	}

}
