package joystick

import (
	"strings"
	"testing"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestDriver() *Driver {
	a := NewAdaptor()
	a.connect = func(j *Adaptor) (err error) {
		j.joystick = &testJoystick{}
		return nil
	}
	a.Connect()
	d := NewDriver(a, "./configs/xbox360_power_a_mini_proex.json")
	d.poll = func() sdl.Event {
		return new(interface{})
	}
	return d
}

func TestJoystickDriverName(t *testing.T) {
	d := initTestDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Joystick"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestDriverStart(t *testing.T) {
	d := initTestDriver()
	d.interval = 1 * time.Millisecond
	gobottest.Assert(t, d.Start(), nil)
	time.Sleep(2 * time.Millisecond)
}

func TestDriverHalt(t *testing.T) {
	d := initTestDriver()
	go func() {
		<-d.halt
	}()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestDriverHandleEvent(t *testing.T) {
	sem := make(chan bool)
	d := initTestDriver()
	d.Start()

	// left x stick
	d.On(d.Event("left_x"), func(data interface{}) {
		gobottest.Assert(t, int16(100), data.(int16))
		sem <- true
	})
	d.handleEvent(&sdl.JoyAxisEvent{
		Which: 0,
		Axis:  0,
		Value: 100,
	})
	select {
	case <-sem:
	case <-time.After(10 * time.Second):
		t.Errorf("Button Event \"left_x\" was not published")
	}

	// x button press
	d.On(d.Event("x_press"), func(data interface{}) {
		sem <- true
	})
	d.handleEvent(&sdl.JoyButtonEvent{
		Which:  0,
		Button: 2,
		State:  1,
	})
	select {
	case <-sem:
	case <-time.After(10 * time.Second):
		t.Errorf("Button Event \"x_press\" was not published")
	}

	// x button  release
	d.On(d.Event("x_release"), func(data interface{}) {
		sem <- true
	})
	d.handleEvent(&sdl.JoyButtonEvent{
		Which:  0,
		Button: 2,
		State:  0,
	})
	select {
	case <-sem:
	case <-time.After(10 * time.Second):
		t.Errorf("Button Event \"x_release\" was not published")
	}

	// down button press
	d.On(d.Event("down"), func(data interface{}) {
		sem <- true
	})
	d.handleEvent(&sdl.JoyHatEvent{
		Which: 0,
		Hat:   0,
		Value: 4,
	})
	select {
	case <-sem:
	case <-time.After(10 * time.Second):
		t.Errorf("Hat Event \"down\" was not published")
	}

	err := d.handleEvent(&sdl.JoyHatEvent{
		Which: 0,
		Hat:   99,
		Value: 4,
	})

	gobottest.Assert(t, err.Error(), "Unknown Hat: 99 4")

	err = d.handleEvent(&sdl.JoyAxisEvent{
		Which: 0,
		Axis:  99,
		Value: 100,
	})

	gobottest.Assert(t, err.Error(), "Unknown Axis: 99")

	err = d.handleEvent(&sdl.JoyButtonEvent{
		Which:  0,
		Button: 99,
		State:  0,
	})

	gobottest.Assert(t, err.Error(), "Unknown Button: 99")
}
