package joystick

import (
	"testing"
	"time"

	"github.com/hybridgroup/go-sdl2/sdl"
	"github.com/hybridgroup/gobot"
)

func initTestJoystickDriver() *JoystickDriver {
	a := NewJoystickAdaptor("bot")
	a.connect = func(j *JoystickAdaptor) (err error) {
		j.joystick = &testJoystick{}
		return nil
	}
	a.Connect()
	d := NewJoystickDriver(a, "bot", "./configs/xbox360_power_a_mini_proex.json")
	d.poll = func() sdl.Event {
		return new(interface{})
	}
	return d
}

func TestJoystickDriverStart(t *testing.T) {
	d := initTestJoystickDriver()
	d.interval = 1 * time.Millisecond
	gobot.Assert(t, len(d.Start()), 0)
	<-time.After(2 * time.Millisecond)
}

func TestJoystickDriverHalt(t *testing.T) {
	d := initTestJoystickDriver()
	go func() {
		<-d.halt
	}()
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestJoystickDriverHandleEvent(t *testing.T) {
	sem := make(chan bool)
	d := initTestJoystickDriver()
	d.Start()

	// left x stick
	gobot.On(d.Event("left_x"), func(data interface{}) {
		gobot.Assert(t, int16(100), data.(int16))
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
	gobot.On(d.Event("x_press"), func(data interface{}) {
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
	gobot.On(d.Event("x_release"), func(data interface{}) {
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
	gobot.On(d.Event("down"), func(data interface{}) {
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

	gobot.Assert(t, err.Error(), "Unknown Hat: 99 4")

	err = d.handleEvent(&sdl.JoyAxisEvent{
		Which: 0,
		Axis:  99,
		Value: 100,
	})

	gobot.Assert(t, err.Error(), "Unknown Axis: 99")

	err = d.handleEvent(&sdl.JoyButtonEvent{
		Which:  0,
		Button: 99,
		State:  0,
	})

	gobot.Assert(t, err.Error(), "Unknown Button: 99")
}
