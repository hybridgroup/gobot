package gpio

import (
	"errors"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MakeyButtonDriver)(nil)

const makeyTestDelay = 250

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
	case <-time.After(makeyTestDelay * time.Millisecond):
		t.Errorf("MakeyButton was not halted")
	}
}

func TestMakeyButtonDriver(t *testing.T) {
	d := initTestMakeyButtonDriver()
	gobottest.Assert(t, d.Pin(), "1")
	gobottest.Refute(t, d.Connection(), nil)
	gobottest.Assert(t, d.interval, 10*time.Millisecond)

	d = NewMakeyButtonDriver(newGpioTestAdaptor(), "1", 30*time.Second)
	gobottest.Assert(t, d.interval, 30*time.Second)
}

func TestMakeyButtonDriverStart(t *testing.T) {
	sem := make(chan bool)
	a := newGpioTestAdaptor()
	d := NewMakeyButtonDriver(a, "1")

	gobottest.Assert(t, d.Start(), nil)

	d.Once(ButtonPush, func(data interface{}) {
		gobottest.Assert(t, d.Active, true)
		sem <- true
	})

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		val = 0
		return
	})

	select {
	case <-sem:
	case <-time.After(makeyTestDelay * time.Millisecond):
		t.Errorf("MakeyButton Event \"Push\" was not published")
	}

	d.Once(ButtonRelease, func(data interface{}) {
		gobottest.Assert(t, d.Active, false)
		sem <- true
	})

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		val = 1
		return
	})

	select {
	case <-sem:
	case <-time.After(makeyTestDelay * time.Millisecond):
		t.Errorf("MakeyButton Event \"Release\" was not published")
	}

	d.Once(Error, func(data interface{}) {
		gobottest.Assert(t, data.(error).Error(), "digital read error")
		sem <- true
	})

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		err = errors.New("digital read error")
		return
	})

	select {
	case <-sem:
	case <-time.After(makeyTestDelay * time.Millisecond):
		t.Errorf("MakeyButton Event \"Error\" was not published")
	}

	// send a halt message
	d.Once(ButtonRelease, func(data interface{}) {
		sem <- true
	})

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		val = 1
		return
	})

	d.halt <- true

	select {
	case <-sem:
		t.Errorf("MakeyButton Event should not have been published")
	case <-time.After(makeyTestDelay * time.Millisecond):
	}
}
