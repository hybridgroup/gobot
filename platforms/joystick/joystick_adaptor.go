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

// JoystickAdaptor represents a connection to a joystick
type JoystickAdaptor struct {
	name     string
	joystick joystick
	connect  func(*JoystickAdaptor) (err error)
}

// NewJoystickAdaptor returns a new JoystickAdaptor with specified name.
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

// Name returns the JoystickAdaptors name
func (j *JoystickAdaptor) Name() string { return j.name }

// Connect connects to the joystick
func (j *JoystickAdaptor) Connect() (errs []error) {
	if err := j.connect(j); err != nil {
		return []error{err}
	}
	return
}

// Finalize closes connection to joystick
func (j *JoystickAdaptor) Finalize() (errs []error) {
	j.joystick.Close()
	return
}
