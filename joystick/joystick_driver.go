package gobotJoystick

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
	JoystickAdaptor *JoystickAdaptor
	config          joystickConfig
}

type pair struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type hat struct {
	Hat  int    `json:"hat"`
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type joystickConfig struct {
	Name    string `json:"name"`
	Guid    string `json:"guid"`
	Axis    []pair `json:"axis"`
	Buttons []pair `json:"buttons"`
	Hats    []hat  `json:"Hats"`
}

type JoystickInterface interface {
}

func NewJoystick(adaptor *JoystickAdaptor) *JoystickDriver {
	d := new(JoystickDriver)
	d.Events = make(map[string]chan interface{})
	d.JoystickAdaptor = adaptor
	d.Commands = []string{}

	var configFile string
	if value, ok := d.JoystickAdaptor.Params["config"]; ok {
		configFile = value.(string)
	} else {
		panic("No joystick config specified")
	}

	file, e := ioutil.ReadFile(configFile)
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

func (me *JoystickDriver) Start() bool {
	go func() {
		var event sdl.Event
		for {
			for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch data := event.(type) {
				case *sdl.JoyAxisEvent:
					if data.Which == me.JoystickAdaptor.joystick.InstanceID() {
						axis := me.findName(data.Axis, me.config.Axis)
						if axis == "" {
							fmt.Println("Unknown Axis:", data.Axis)
						} else {
							gobot.Publish(me.Events[axis], data.Value)
						}
					}
				case *sdl.JoyButtonEvent:
					if data.Which == me.JoystickAdaptor.joystick.InstanceID() {
						button := me.findName(data.Button, me.config.Buttons)
						if button == "" {
							fmt.Println("Unknown Button:", data.Button)
						} else {
							if data.State == 1 {
								gobot.Publish(me.Events[fmt.Sprintf("%s_press", button)], nil)
							} else {
								gobot.Publish(me.Events[fmt.Sprintf("%s_release", button)], nil)
							}
						}
					}
				case *sdl.JoyHatEvent:
					if data.Which == me.JoystickAdaptor.joystick.InstanceID() {
						hat := me.findHatName(data.Value, data.Hat, me.config.Hats)
						if hat == "" {
							fmt.Println("Unknown Hat:", data.Hat, data.Value)
						} else {
							gobot.Publish(me.Events[hat], true)
						}
					}
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	return true
}
func (me *JoystickDriver) Init() bool { return true }
func (me *JoystickDriver) Halt() bool { return true }

func (me *JoystickDriver) findName(id uint8, list []pair) string {
	for _, value := range list {
		if int(id) == value.Id {
			return value.Name
		}
	}
	return ""
}

func (me *JoystickDriver) findHatName(id uint8, hat uint8, list []hat) string {
	for _, value := range list {
		if int(id) == value.Id && int(hat) == value.Hat {
			return value.Name
		}
	}
	return ""
}
