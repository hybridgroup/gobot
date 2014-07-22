package joystick

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hybridgroup/go-sdl2/sdl"
	"github.com/hybridgroup/gobot"
)

type JoystickDriver struct {
	gobot.Driver
	config joystickConfig
	poll   func() sdl.Event
}

type pair struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type hat struct {
	Hat  int    `json:"hat"`
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type joystickConfig struct {
	Name    string `json:"name"`
	Guid    string `json:"guid"`
	Axis    []pair `json:"axis"`
	Buttons []pair `json:"buttons"`
	Hats    []hat  `json:"Hats"`
}

func NewJoystickDriver(a *JoystickAdaptor, name string, config string) *JoystickDriver {
	d := &JoystickDriver{
		Driver: *gobot.NewDriver(
			name,
			"JoystickDriver",
			a,
		),
		poll: func() sdl.Event {
			return sdl.PollEvent()
		},
	}

	file, e := ioutil.ReadFile(config)
	if e != nil {
		panic(fmt.Sprintf("File error: %v\n", e))
	}
	var jsontype joystickConfig
	json.Unmarshal(file, &jsontype)
	d.config = jsontype
	for _, value := range d.config.Buttons {
		d.AddEvent(fmt.Sprintf("%s_press", value.Name))
		d.AddEvent(fmt.Sprintf("%s_release", value.Name))
	}
	for _, value := range d.config.Axis {
		d.AddEvent(value.Name)
	}
	for _, value := range d.config.Hats {
		d.AddEvent(value.Name)
	}
	return d
}

func (j *JoystickDriver) adaptor() *JoystickAdaptor {
	return j.Adaptor().(*JoystickAdaptor)
}

func (j *JoystickDriver) Start() bool {
	gobot.Every(j.Interval(), func() {
		event := j.poll()
		if event != nil {
			j.handleEvent(event)
		}
	})
	return true
}

func (j *JoystickDriver) handleEvent(event sdl.Event) error {
	switch data := event.(type) {
	case *sdl.JoyAxisEvent:
		if data.Which == j.adaptor().joystick.InstanceID() {
			axis := j.findName(data.Axis, j.config.Axis)
			if axis == "" {
				e := errors.New(fmt.Sprintf("Unknown Axis: %v", data.Axis))
				fmt.Println(e.Error())
				return e
			} else {
				gobot.Publish(j.Event(axis), data.Value)
			}
		}
	case *sdl.JoyButtonEvent:
		if data.Which == j.adaptor().joystick.InstanceID() {
			button := j.findName(data.Button, j.config.Buttons)
			if button == "" {
				e := errors.New(fmt.Sprintf("Unknown Button: %v", data.Button))
				fmt.Println(e.Error())
				return e
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
				e := errors.New(fmt.Sprintf("Unknown Hat: %v %v", data.Hat, data.Value))
				fmt.Println(e.Error())
				return e
			} else {
				gobot.Publish(j.Event(hat), true)
			}
		}
	}
	return nil
}

func (j *JoystickDriver) Halt() bool { return true }

func (j *JoystickDriver) findName(id uint8, list []pair) string {
	for _, value := range list {
		if int(id) == value.ID {
			return value.Name
		}
	}
	return ""
}

func (j *JoystickDriver) findHatName(id uint8, hat uint8, list []hat) string {
	for _, value := range list {
		if int(id) == value.ID && int(hat) == value.Hat {
			return value.Name
		}
	}
	return ""
}
