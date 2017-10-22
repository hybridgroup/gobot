package gpio

import (
	"errors"
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*ButtonDriver)(nil)

const buttonTestDelay = 250

func initTestButtonDriver() *ButtonDriver {
	return NewButtonDriver(newGpioTestAdaptor(), "1")
}

func TestButtonDriverHalt(t *testing.T) {
	d := initTestButtonDriver()
	go func() {
		<-d.halt
	}()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestButtonDriver(t *testing.T) {
	d := NewButtonDriver(newGpioTestAdaptor(), "1")
	gobottest.Refute(t, d.Connection(), nil)

	d = NewButtonDriver(newGpioTestAdaptor(), "1", 30*time.Second)
	gobottest.Assert(t, d.interval, 30*time.Second)
}

func TestButtonDriverStart(t *testing.T) {
	sem := make(chan bool, 0)
	a := newGpioTestAdaptor()
	d := NewButtonDriver(a, "1")

	d.Once(ButtonPush, func(data interface{}) {
		gobottest.Assert(t, d.Active, true)
		sem <- true
	})

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		val = 1
		return
	})

	gobottest.Assert(t, d.Start(), nil)

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Push\" was not published")
	}

	d.Once(ButtonRelease, func(data interface{}) {
		gobottest.Assert(t, d.Active, false)
		sem <- true
	})

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		val = 0
		return
	})

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Release\" was not published")
	}

	d.Once(Error, func(data interface{}) {
		sem <- true
	})

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		err = errors.New("digital read error")
		return
	})

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Error\" was not published")
	}

	d.Once(ButtonPush, func(data interface{}) {
		sem <- true
	})

	d.halt <- true

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		val = 1
		return
	})

	select {
	case <-sem:
		t.Errorf("Button Event \"Press\" should not published")
	case <-time.After(buttonTestDelay * time.Millisecond):
	}
}

func TestButtonDriverDefaultState(t *testing.T) {
	sem := make(chan bool, 0)
	a := newGpioTestAdaptor()
	d := NewButtonDriver(a, "1")
	d.DefaultState = 1

	d.Once(ButtonPush, func(data interface{}) {
		gobottest.Assert(t, d.Active, true)
		sem <- true
	})

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		val = 0
		return
	})

	gobottest.Assert(t, d.Start(), nil)

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Push\" was not published")
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
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Release\" was not published")
	}
}

func TestButtonDriverDefaultName(t *testing.T) {
	g := initTestButtonDriver()
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "Button"), true)
}

func TestButtonDriverSetName(t *testing.T) {
	g := initTestButtonDriver()
	g.SetName("mybot")
	gobottest.Assert(t, g.Name(), "mybot")
}
