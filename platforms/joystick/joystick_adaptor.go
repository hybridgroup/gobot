package joystick

import (
	"errors"
	"github.com/hybridgroup/go-sdl2/sdl"
	"github.com/hybridgroup/gobot"
)

var _ gobot.AdaptorInterface = (*JoystickAdaptor)(nil)

type joystick interface {
	Close()
	InstanceID() sdl.JoystickID
}

type JoystickAdaptor struct {
	gobot.Adaptor
	joystick joystick
	connect  func(*JoystickAdaptor) (err error)
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
		connect: func(j *JoystickAdaptor) (err error) {
			sdl.Init(sdl.INIT_JOYSTICK)
			if sdl.NumJoysticks() > 0 {
				j.joystick = sdl.JoystickOpen(0)
				return
			}
			return errors.New("No joystick available")
		},
	}
}

// Connect returns true if connection to device is succesfull
func (j *JoystickAdaptor) Connect() error {
	return j.connect(j)
}

// Finalize closes connection to device
func (j *JoystickAdaptor) Finalize() error {
	j.joystick.Close()
	return nil
}
