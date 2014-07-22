package joystick

import "github.com/hybridgroup/go-sdl2/sdl"

type testJoystick struct{}

func (t *testJoystick) Close()                     {}
func (t *testJoystick) InstanceID() sdl.JoystickID { return 0 }
