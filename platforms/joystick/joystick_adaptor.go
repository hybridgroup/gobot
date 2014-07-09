package joystick

import (
	"github.com/hybridgroup/go-sdl2/sdl"
	"github.com/hybridgroup/gobot"
)

type JoystickAdaptor struct {
	gobot.Adaptor
	joystick *sdl.Joystick
	connect  func(*JoystickAdaptor)
}

func NewJoystickAdaptor(name string) *JoystickAdaptor {
	return &JoystickAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"JoystickAdaptor",
		),
		connect: func(j *JoystickAdaptor) {
			sdl.Init(sdl.INIT_JOYSTICK)
			if sdl.NumJoysticks() > 0 {
				j.joystick = sdl.JoystickOpen(0)
			} else {
				panic("No joystick available")
			}
		},
	}
}

func (j *JoystickAdaptor) Connect() bool {
	j.connect(j)
	return true
}

func (j *JoystickAdaptor) Finalize() bool {
	j.joystick.Close()
	return true
}
