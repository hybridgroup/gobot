package joystick

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hybridgroup/go-sdl2/sdl"
	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*JoystickDriver)(nil)

type JoystickDriver struct {
	name       string
	interval   time.Duration
	connection gobot.Connection
	configPath string
	config     joystickConfig
	poll       func() sdl.Event
	gobot.Eventer
}

// pair is a JSON representation of name and id
type pair struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// hat is a JSON representation of hat, name and id
type hat struct {
	Hat  int    `json:"hat"`
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// joystickConfig is a JSON representation of configuration values
type joystickConfig struct {
	Name    string `json:"name"`
	Guid    string `json:"guid"`
	Axis    []pair `json:"axis"`
	Buttons []pair `json:"buttons"`
	Hats    []hat  `json:"Hats"`
}

// NewJoystickDriver creates a joyscript driver by name.
//
// It adds the following events:
//     (button)_press - triggered when (button) is pressed
//     (button)_release - triggered when (button) is released
func NewJoystickDriver(a *JoystickAdaptor, name string, config string, v ...time.Duration) *JoystickDriver {
	d := &JoystickDriver{
		name:       name,
		connection: a,
		Eventer:    gobot.NewEventer(),
		configPath: config,
		poll: func() sdl.Event {
			return sdl.PollEvent()
		},
		interval: 10 * time.Millisecond,
	}

	if len(v) > 0 {
		d.interval = v[0]
	}

	d.AddEvent("error")
	return d
}
func (j *JoystickDriver) Name() string                 { return j.name }
func (j *JoystickDriver) Connection() gobot.Connection { return j.connection }

// adaptor returns joystick adaptor
func (j *JoystickDriver) adaptor() *JoystickAdaptor {
	return j.Connection().(*JoystickAdaptor)
}

// Start initiallizes event polling with defined interval
func (j *JoystickDriver) Start() (errs []error) {
	file, err := ioutil.ReadFile(j.configPath)
	if err != nil {
		return []error{err}
	}

	var jsontype joystickConfig
	json.Unmarshal(file, &jsontype)
	j.config = jsontype

	for _, value := range j.config.Buttons {
		j.AddEvent(fmt.Sprintf("%s_press", value.Name))
		j.AddEvent(fmt.Sprintf("%s_release", value.Name))
	}
	for _, value := range j.config.Axis {
		j.AddEvent(value.Name)
	}
	for _, value := range j.config.Hats {
		j.AddEvent(value.Name)
	}

	go func() {
		for {
			event := j.poll()
			if event != nil {
				if err = j.handleEvent(event); err != nil {
					gobot.Publish(j.Event("error"), err)
				}
			}
			<-time.After(j.interval)
		}
	}()
	return
}

// Halt stops joystick driver
func (j *JoystickDriver) Halt() (errs []error) { return }

// HandleEvent publishes an specific event according to data received
func (j *JoystickDriver) handleEvent(event sdl.Event) error {
	switch data := event.(type) {
	case *sdl.JoyAxisEvent:
		if data.Which == j.adaptor().joystick.InstanceID() {
			axis := j.findName(data.Axis, j.config.Axis)
			if axis == "" {
				return errors.New(fmt.Sprintf("Unknown Axis: %v", data.Axis))
			} else {
				gobot.Publish(j.Event(axis), data.Value)
			}
		}
	case *sdl.JoyButtonEvent:
		if data.Which == j.adaptor().joystick.InstanceID() {
			button := j.findName(data.Button, j.config.Buttons)
			if button == "" {
				return errors.New(fmt.Sprintf("Unknown Button: %v", data.Button))
			} else {
				if data.State == 1 {
					gobot.Publish(j.Event(fmt.Sprintf("%s_press", button)), nil)
				} else {
					gobot.Publish(j.Event(fmt.Sprintf("%s_release", button)), nil)
				}
			}
		}
	case *sdl.JoyHatEvent:
		if data.Which == j.adaptor().joystick.InstanceID() {
			hat := j.findHatName(data.Value, data.Hat, j.config.Hats)
			if hat == "" {
				return errors.New(fmt.Sprintf("Unknown Hat: %v %v", data.Hat, data.Value))
			} else {
				gobot.Publish(j.Event(hat), true)
			}
		}
	}
	return nil
}

func (j *JoystickDriver) findName(id uint8, list []pair) string {
	for _, value := range list {
		if int(id) == value.ID {
			return value.Name
		}
	}
	return ""
}

// findHatName returns name from hat found by id in provided list
func (j *JoystickDriver) findHatName(id uint8, hat uint8, list []hat) string {
	for _, lHat := range list {
		if int(id) == lHat.ID && int(hat) == lHat.Hat {
			return lHat.Name
		}
	}
	return ""
}
