package joystick

import (
	"errors"

	"github.com/hybridgroup/go-sdl2/sdl"
	"github.com/hybridgroup/gobot"
)

var _ gobot.Adaptor = (*JoystickAdaptor)(nil)

type joystick interface {
	Close()
	InstanceID() sdl.JoystickID
}

type JoystickAdaptor struct {
	name     string
	joystick joystick
	connect  func(*JoystickAdaptor) (err error)
}

// NewJoysctickAdaptor creates a new adaptor with specified name.
// It creates a connect function to joystick in position 0
// or panics if no joystick can be found
func NewJoystickAdaptor(name string) *JoystickAdaptor {
	return &JoystickAdaptor{
		name: name,
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

func (j *JoystickAdaptor) Name() string { return j.name }

// Connect returns true if connection to device is succesfull
func (j *JoystickAdaptor) Connect() (errs []error) {
	if err := j.connect(j); err != nil {
		return []error{err}
	}
	return
}

// Finalize closes connection to device
func (j *JoystickAdaptor) Finalize() (errs []error) {
	j.joystick.Close()
	return
}
