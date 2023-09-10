package joystick

import js "github.com/0xcafed00d/joystick"

type testJoystick struct {
	axisCount   int
	buttonCount int
}

func (t *testJoystick) Close()                  {}
func (t *testJoystick) ID() int                 { return 0 }
func (t *testJoystick) ButtonCount() int        { return t.buttonCount }
func (t *testJoystick) AxisCount() int          { return t.axisCount }
func (t *testJoystick) Name() string            { return "test-joy" }
func (t *testJoystick) Read() (js.State, error) { return js.State{}, nil }
