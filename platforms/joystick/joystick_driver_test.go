//nolint:forcetypeassert // ok here
package joystick

import (
	"strings"
	"testing"
	"time"

	js "github.com/0xcafed00d/joystick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestDriver(config string) (*Driver, *testJoystick) {
	a := NewAdaptor("6")
	tj := &testJoystick{}
	a.connect = func(j *Adaptor) error {
		j.joystick = tj
		return nil
	}
	_ = a.Connect()
	d := NewDriver(a, config)
	return d, tj
}

func TestJoystickDriverName(t *testing.T) {
	d, _ := initTestDriver("./configs/dualshock3.json")
	assert.True(t, strings.HasPrefix(d.Name(), "Joystick"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestDriverStart(t *testing.T) {
	d, _ := initTestDriver("./configs/dualshock3.json")
	d.interval = 1 * time.Millisecond
	require.NoError(t, d.Start())
	time.Sleep(2 * time.Millisecond)
}

func TestDriverHalt(t *testing.T) {
	d, _ := initTestDriver("./configs/dualshock3.json")
	go func() {
		<-d.halt
	}()
	require.NoError(t, d.Halt())
}

func TestDriverHandleEventDS3(t *testing.T) {
	sem := make(chan bool)
	d, tj := initTestDriver("dualshock3")
	tj.axisCount = 6
	tj.buttonCount = 17

	if err := d.initConfig(); err != nil {
		require.Fail(t, "initConfig() error: %v", err)
	}

	d.initEvents()

	// left x stick
	_ = d.On(d.Event("left_x"), func(data interface{}) {
		assert.Equal(t, int(255), data.(int))
		sem <- true
	})
	_ = d.handleAxes(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  0,
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		require.Fail(t, "Button Event \"left_x\" was not published")
	}

	// square button press
	_ = d.On(d.Event("square_press"), func(data interface{}) {
		sem <- true
	})
	_ = d.handleButtons(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  1 << d.findID("square", d.config.Buttons),
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		require.Fail(t, "Button Event \"square_press\" was not published")
	}

	// square button release
	_ = d.On(d.Event("square_release"), func(data interface{}) {
		sem <- true
	})
	_ = d.handleButtons(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  0,
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		require.Fail(t, "Button Event \"square_release\" was not published")
	}
}

func TestDriverHandleEventJSONDS3(t *testing.T) {
	sem := make(chan bool)
	d, tj := initTestDriver("./configs/dualshock3.json")
	tj.axisCount = 6
	tj.buttonCount = 17

	if err := d.initConfig(); err != nil {
		require.Fail(t, "initConfig() error: %v", err)
	}

	d.initEvents()

	// left x stick
	_ = d.On(d.Event("left_x"), func(data interface{}) {
		assert.Equal(t, int(255), data.(int))
		sem <- true
	})
	_ = d.handleAxes(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  0,
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		require.Fail(t, "Button Event \"left_x\" was not published")
	}

	// square button press
	_ = d.On(d.Event("square_press"), func(data interface{}) {
		sem <- true
	})
	_ = d.handleButtons(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  1 << d.findID("square", d.config.Buttons),
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		require.Fail(t, "Button Event \"square_press\" was not published")
	}

	// square button release
	_ = d.On(d.Event("square_release"), func(data interface{}) {
		sem <- true
	})
	_ = d.handleButtons(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  0,
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		require.Fail(t, "Button Event \"square_release\" was not published")
	}
}

func TestDriverHandleEventDS4(t *testing.T) {
	sem := make(chan bool)
	d, tj := initTestDriver("dualshock4")
	tj.axisCount = 6
	tj.buttonCount = 17

	if err := d.initConfig(); err != nil {
		require.Fail(t, "initConfig() error: %v", err)
	}

	d.initEvents()

	// left x stick
	_ = d.On(d.Event("left_x"), func(data interface{}) {
		assert.Equal(t, int(255), data.(int))
		sem <- true
	})
	_ = d.handleAxes(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  0,
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		require.Fail(t, "Button Event \"left_x\" was not published")
	}

	// square button press
	_ = d.On(d.Event("square_press"), func(data interface{}) {
		sem <- true
	})
	_ = d.handleButtons(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  1 << d.findID("square", d.config.Buttons),
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		require.Fail(t, "Button Event \"square_press\" was not published")
	}

	// square button release
	_ = d.On(d.Event("square_release"), func(data interface{}) {
		sem <- true
	})
	_ = d.handleButtons(js.State{
		AxisData: []int{255, 0, 0, 0, 0, 0},
		Buttons:  0,
	})
	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		require.Fail(t, "Button Event \"square_release\" was not published")
	}
}

func TestDriverInvalidConfig(t *testing.T) {
	d, _ := initTestDriver("./configs/doesnotexist")
	err := d.Start()
	require.ErrorContains(t, err, "loadfile error")
}
