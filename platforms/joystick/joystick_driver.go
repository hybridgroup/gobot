package joystick

import (
	"encoding/json"
	"fmt"
	"github.com/hybridgroup/go-sdl2/sdl"
	"github.com/hybridgroup/gobot"
	"io/ioutil"
	"time"
)

type JoystickDriver struct {
	gobot.Driver
	Adaptor *JoystickAdaptor
	config  joystickConfig
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
		Driver: gobot.Driver{
			Events: make(map[string]chan interface{}),
		},
		Adaptor: a,
	}

	file, e := ioutil.ReadFile(config)
	if e != nil {
		panic(fmt.Sprintf("File error: %v\n", e))
	}
	var jsontype joystickConfig
	json.Unmarshal(file, &jsontype)
	d.config = jsontype
	for _, value := range d.config.Buttons {
		d.Events[fmt.Sprintf("%s_press", value.Name)] = make(chan interface{}, 0)
		d.Events[fmt.Sprintf("%s_release", value.Name)] = make(chan interface{}, 0)
	}
	for _, value := range d.config.Axis {
		d.Events[value.Name] = make(chan interface{}, 0)
	}
	for _, value := range d.config.Hats {
		d.Events[value.Name] = make(chan interface{}, 0)
	}
	return d
}

func (j *JoystickDriver) Start() bool {
	go func() {
		var event sdl.Event
		for {
			for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch data := event.(type) {
				case *sdl.JoyAxisEvent:
					if data.Which == j.Adaptor.joystick.InstanceID() {
						axis := j.findName(data.Axis, j.config.Axis)
						if axis == "" {
							fmt.Println("Unknown Axis:", data.Axis)
						} else {
							gobot.Publish(j.Events[axis], data.Value)
						}
					}
				case *sdl.JoyButtonEvent:
					if data.Which == j.Adaptor.joystick.InstanceID() {
						button := j.findName(data.Button, j.config.Buttons)
						if button == "" {
							fmt.Println("Unknown Button:", data.Button)
						} else {
							if data.State == 1 {
								gobot.Publish(j.Events[fmt.Sprintf("%s_press", button)], nil)
							} else {
								gobot.Publish(j.Events[fmt.Sprintf("%s_release", button)], nil)
							}
						}
					}
				case *sdl.JoyHatEvent:
					if data.Which == j.Adaptor.joystick.InstanceID() {
						hat := j.findHatName(data.Value, data.Hat, j.config.Hats)
						if hat == "" {
							fmt.Println("Unknown Hat:", data.Hat, data.Value)
						} else {
							gobot.Publish(j.Events[hat], true)
						}
					}
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	return true
}
func (j *JoystickDriver) Init() bool { return true }
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
