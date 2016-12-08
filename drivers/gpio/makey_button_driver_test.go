package gpio

import (
	"errors"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MakeyButtonDriver)(nil)

const MAKEY_TEST_DELAY = 30

func initTestMakeyButtonDriver() *MakeyButtonDriver {
	return NewMakeyButtonDriver(newGpioTestAdaptor(), "1")
}

func TestMakeyButtonDriverHalt(t *testing.T) {
	d := initTestMakeyButtonDriver()
	done := make(chan struct{})
	go func() {
		<-d.halt
		close(done)
	}()
	gobottest.Assert(t, d.Halt(), nil)
	select {
	case <-done:
	case <-time.After(time.Millisecond):
		t.Errorf("MakeyButton was not halted")
	}
}

func TestMakeyButtonDriver(t *testing.T) {
	d := initTestMakeyButtonDriver()
	gobottest.Assert(t, d.Pin(), "1")
	gobottest.Refute(t, d.Connection(), nil)
	gobottest.Assert(t, d.interval, 10*time.Millisecond)

	d = NewMakeyButtonDriver(newGpioTestAdaptor(), "1", 30*time.Second)
	gobottest.Assert(t, d.interval, MAKEY_TEST_DELAY*time.Second)
}

func TestMakeyButtonDriverStart(t *testing.T) {
	sem := make(chan bool)
	d := initTestMakeyButtonDriver()
	gobottest.Assert(t, d.Start(), nil)

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 0
		return
	}

	d.Once(ButtonPush, func(data interface{}) {
		gobottest.Assert(t, d.Active, true)
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(MAKEY_TEST_DELAY * time.Millisecond):
		t.Errorf("MakeyButton Event \"Push\" was not published")
	}

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}

	d.Once(ButtonRelease, func(data interface{}) {
		gobottest.Assert(t, d.Active, false)
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(MAKEY_TEST_DELAY * time.Millisecond):
		t.Errorf("MakeyButton Event \"Release\" was not published")
	}

	testAdaptorDigitalRead = func() (val int, err error) {
		err = errors.New("digital read error")
		return
	}

	d.Once(Error, func(data interface{}) {
		gobottest.Assert(t, data.(error).Error(), "digital read error")
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(MAKEY_TEST_DELAY * time.Millisecond):
		t.Errorf("MakeyButton Event \"Error\" was not published")
	}

	// send a halt message
	d.Once(ButtonRelease, func(data interface{}) {
		sem <- true
	})

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}

	d.halt <- true

	select {
	case <-sem:
		t.Errorf("MakeyButton Event should not have been published")
	case <-time.After(MAKEY_TEST_DELAY * time.Millisecond):
	}
}
