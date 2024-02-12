package joystick

import (
	"fmt"
	"strconv"

	js "github.com/0xcafed00d/joystick"

	"gobot.io/x/gobot/v2"
)

// Adaptor represents a connection to a joystick
type Adaptor struct {
	name     string
	id       string
	joystick js.Joystick
	connect  func(*Adaptor) error
}

// NewAdaptor returns a new Joystick Adaptor.
// Pass in the ID of the joystick you wish to connect to.
func NewAdaptor(id string) *Adaptor {
	return &Adaptor{
		name: gobot.DefaultName("Joystick"),
		connect: func(j *Adaptor) error {
			i, err := strconv.Atoi(id)
			if err != nil {
				return fmt.Errorf("invalid joystick ID: %v", err)
			}

			joy, err := js.Open(i)
			if err != nil {
				return fmt.Errorf("no joystick available: %v", err)
			}

			j.id = id
			j.joystick = joy
			return nil
		},
	}
}

// Name returns the adaptors name
func (j *Adaptor) Name() string { return j.name }

// SetName sets the adaptors name
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
