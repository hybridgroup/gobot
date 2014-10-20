package joystick

import (
	"github.com/hybridgroup/go-sdl2/sdl"
	"github.com/hybridgroup/gobot"
)

type joystick interface {
	Close()
	InstanceID() sdl.JoystickID
}

type JoystickAdaptor struct {
	gobot.Adaptor
	joystick joystick
	connect  func(*JoystickAdaptor)
}

// NewJoysctickAdaptor creates a new adaptor with specified name.
// It creates a connect function to joystick in position 0
// or panics if no joystick can be found
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

// Connect returns true if connection to device is succesfull
func (j *JoystickAdaptor) Connect() bool {
	j.connect(j)
	return true
}

// Finalize closes connection to device
func (j *JoystickAdaptor) Finalize() bool {
	j.joystick.Close()
	return true
}
