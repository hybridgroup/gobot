package gpio

import (
	"errors"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

func initTestMakeyButtonDriver() *MakeyButtonDriver {
	return NewMakeyButtonDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestMakeyButtonDriverHalt(t *testing.T) {
	d := initTestMakeyButtonDriver()
	go func() {
		<-d.halt
	}()
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestMakeyButtonDriver(t *testing.T) {
	d := NewMakeyButtonDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
	gobot.Assert(t, d.Name(), "bot")
	gobot.Assert(t, d.Connection().Name(), "adaptor")

	d = NewMakeyButtonDriver(newGpioTestAdaptor("adaptor"), "bot", "1", 30*time.Second)
	gobot.Assert(t, d.interval, 30*time.Second)
}

func TestMakeyButtonDriverStart(t *testing.T) {
	sem := make(chan bool, 0)
	d := initTestMakeyButtonDriver()
	gobot.Assert(t, len(d.Start()), 0)

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 0
		return
	}

	gobot.Once(d.Event(Push), func(data interface{}) {
		gobot.Assert(t, d.Active, true)
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(15 * time.Millisecond):
		t.Errorf("MakeyButton Event \"Push\" was not published")
	}

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}

	gobot.Once(d.Event(Release), func(data interface{}) {
		gobot.Assert(t, d.Active, false)
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(15 * time.Millisecond):
		t.Errorf("MakeyButton Event \"Release\" was not published")
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
		t.Errorf("MakeyButton Event \"Error\" was not published")
	}

	gobot.Once(d.Event(Release), func(data interface{}) {
		sem <- true
	})
	testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}

	d.halt <- true

	select {
	case <-sem:
		t.Errorf("MakeyButton Event should not published")
	case <-time.After(30 * time.Millisecond):
	}
}
