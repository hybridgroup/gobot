package joystick

import (
	"fmt"

	"gobot.io/x/gobot/v2"

	js "github.com/0xcafed00d/joystick"
)

type joystick interface {
	Close()
	ID() int
}

// Adaptor represents a connection to a joystick
type Adaptor struct {
	name     string
	id int
	joystick js.Joystick
	connect  func(*Adaptor) error
}

// NewAdaptor returns a new Joystick Adaptor.
// Pass in the ID of the joystick you wish to connect to.
func NewAdaptor(id int) *Adaptor {
	return &Adaptor{
		name: gobot.DefaultName("Joystick"),
		connect: func(j *Adaptor) error {
			joy, err := js.Open(id)
			if err != nil {
				return fmt.Errorf("No joystick available: %v", err)
			}

			j.id = id
			j.joystick = joy
			return nil
		},
	}
}

// Name returns the Adaptors name
func (j *Adaptor) Name() string { return j.name }

// SetName sets the Adaptors name
func (j *Adaptor) SetName(n string) { j.name = n }

// Connect connects to the joystick
func (j *Adaptor) Connect() error {
	return j.connect(j)
}

// Finalize closes connection to joystick
func (j *Adaptor) Finalize() error {
	j.joystick.Close()
	return nil
}
