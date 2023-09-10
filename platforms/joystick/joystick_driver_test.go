package joystick

import (
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"

	js "github.com/0xcafed00d/joystick"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestDriver(config string) (*Driver, *testJoystick) {
	a := NewAdaptor(6)
	tj := &testJoystick{}
	a.connect = func(j *Adaptor) (err error) {
		j.joystick = tj
		return nil
	}
	_ = a.Connect()
	d := NewDriver(a, config)
	return d, tj
}

func TestJoystickDriverName(t *testing.T) {
	d, _ := initTestDriver("./configs/dualshock3.json")
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Joystick"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestDriverStart(t *testing.T) {
	d, _ := initTestDriver("./configs/dualshock3.json")
	d.interval = 1 * time.Millisecond
	gobottest.Assert(t, d.Start(), nil)
	time.Sleep(2 * time.Millisecond)
}

func TestDriverHalt(t *testing.T) {
	d, _ := initTestDriver("./configs/dualshock3.json")
	go func() {
		<-d.halt
	}()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestDriverHandleEvent(t *testing.T) {
	sem := make(chan bool)
	d, tj := initTestDriver("./configs/dualshock3.json")
	tj.axisCount = 6
	tj.buttonCount = 17

	if err := d.initConfig(); err != nil {
		t.Errorf("initConfig() error: %v", err)
	}
	
	d.initEvents()

	// left x stick
	_ = d.On(d.Event("left_x"), func(data interface{}) {
		gobottest.Assert(t, int(255), data.(int))
		sem <- true
	})
	_ = d.handleAxes(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  0,
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf("Button Event \"left_x\" was not published")
	}

	// x button press
	_ = d.On(d.Event("x_press"), func(data interface{}) {
		sem <- true
	})
	_ = d.handleButtons(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  1 << 14,
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf("Button Event \"x_press\" was not published")
	}

	// x button release
	_ = d.On(d.Event("x_release"), func(data interface{}) {
		sem <- true
	})
	_ = d.handleButtons(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  0,
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf("Button Event \"x_release\" was not published")
	}
}

func TestDriverInvalidConfig(t *testing.T) {
	d, _ := initTestDriver("./configs/doesnotexist")
	err := d.Start()
	gobottest.Assert(t, strings.Contains(err.Error(), "open ./configs/doesnotexist: no such file or directory"), true)
}
