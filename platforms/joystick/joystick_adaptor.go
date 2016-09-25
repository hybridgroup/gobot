package joystick

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
)

type joystick interface {
	Close()
	InstanceID() sdl.JoystickID
}

// Adaptor represents a connection to a joystick
type Adaptor struct {
	name     string
	joystick joystick
	connect  func(*Adaptor) (err error)
}

// NewAdaptor returns a new Joystick Adaptor.
func NewAdaptor() *Adaptor {
	return &Adaptor{
		connect: func(j *Adaptor) (err error) {
			sdl.Init(sdl.INIT_JOYSTICK)
			if sdl.NumJoysticks() > 0 {
				j.joystick = sdl.JoystickOpen(0)
				return
			}
			return errors.New("No joystick available")
		},
	}
}

// Name returns the Adaptors name
func (j *Adaptor) Name() string { return j.name }

// SetName sets the Adaptors name
func (j *Adaptor) SetName(n string) { j.name = n }

// Connect connects to the joystick
func (j *Adaptor) Connect() (errs []error) {
	if err := j.connect(j); err != nil {
		return []error{err}
	}
	return
}

// Finalize closes connection to joystick
func (j *Adaptor) Finalize() (errs []error) {
	j.joystick.Close()
	return
}
