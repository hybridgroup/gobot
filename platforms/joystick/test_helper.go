package joystick

import js "github.com/0xcafed00d/joystick"

type testJoystick struct{}

func (t *testJoystick) Close()                     {}
func (t *testJoystick) ID() int { return 0 }
func (t *testJoystick) ButtonCount() int { return 0 }
func (t *testJoystick) AxisCount() int { return 0 }
func (t *testJoystick) Name() string { return "test-joy" }
func (t *testJoystick) Read() (js.State, error) { return js.State{}, nil }
