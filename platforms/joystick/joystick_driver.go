package joystick

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"gobot.io/x/gobot"
)

// Driver represents a joystick
type Driver struct {
	name       string
	interval   time.Duration
	connection gobot.Connection
	configPath string
	config     joystickConfig
	poll       func() sdl.Event
	halt       chan bool
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
	GUID    string `json:"guid"`
	Axis    []pair `json:"axis"`
	Buttons []pair `json:"buttons"`
	Hats    []hat  `json:"Hats"`
}

// NewDriver returns a new Driver with a polling interval of
// 10 Milliseconds given a Joystick Adaptor and json button configuration
// file location.
//
// Optionally accepts:
//  time.Duration: Interval at which the Driver is polled for new information
func NewDriver(a *Adaptor, config string, v ...time.Duration) *Driver {
	d := &Driver{
		name:       gobot.DefaultName("Joystick"),
		connection: a,
		Eventer:    gobot.NewEventer(),
		configPath: config,
		poll: func() sdl.Event {
			return sdl.PollEvent()
		},
		interval: 10 * time.Millisecond,
		halt:     make(chan bool, 0),
	}

	if len(v) > 0 {
		d.interval = v[0]
	}

	d.AddEvent("error")
	return d
}

// Name returns the Drivers name
func (j *Driver) Name() string { return j.name }

// SetName sets the Drivers name
func (j *Driver) SetName(n string) { j.name = n }

// Connection returns the Drivers connection
func (j *Driver) Connection() gobot.Connection { return j.connection }

// adaptor returns joystick adaptor
func (j *Driver) adaptor() *Adaptor {
	return j.Connection().(*Adaptor)
}

// Start and polls the state of the joystick at the given interval.
//
// Emits the Events:
//	Error error - On button error
//	Events defined in the json button configuration file.
//	They will have the format:
//		[button]_press
//		[button]_release
//		[axis]
func (j *Driver) Start() (err error) {
	file, e := ioutil.ReadFile(j.configPath)
	if e != nil {
		return e
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
			for event := j.poll(); event != nil; event = j.poll() {
				if errs := j.handleEvent(event); errs != nil {
					j.Publish(j.Event("error"), errs)
				}
			}
			select {
			case <-time.After(j.interval):
			case <-j.halt:
				return
			}
		}
	}()
	return
}

// Halt stops joystick driver
func (j *Driver) Halt() (err error) {
	j.halt <- true
	return
}

// HandleEvent publishes an specific event according to data received
func (j *Driver) handleEvent(event sdl.Event) error {
	switch data := event.(type) {
	case *sdl.JoyAxisEvent:
		if data.Which == j.adaptor().joystick.InstanceID() {
			axis := j.findName(data.Axis, j.config.Axis)
			if axis == "" {
				return fmt.Errorf("Unknown Axis: %v", data.Axis)
			}
			j.Publish(j.Event(axis), data.Value)
		}
	case *sdl.JoyButtonEvent:
		if data.Which == j.adaptor().joystick.InstanceID() {
			button := j.findName(data.Button, j.config.Buttons)
			if button == "" {
				return fmt.Errorf("Unknown Button: %v", data.Button)
			}
			if data.State == 1 {
				j.Publish(j.Event(fmt.Sprintf("%s_press", button)), nil)
			}
			j.Publish(j.Event(fmt.Sprintf("%s_release", button)), nil)
		}
	case *sdl.JoyHatEvent:
		if data.Which == j.adaptor().joystick.InstanceID() {
			hat := j.findHatName(data.Value, data.Hat, j.config.Hats)
			if hat == "" {
				return fmt.Errorf("Unknown Hat: %v %v", data.Hat, data.Value)
			}
			j.Publish(j.Event(hat), true)
		}
	}
	return nil
}

func (j *Driver) findName(id uint8, list []pair) string {
	for _, value := range list {
		if int(id) == value.ID {
			return value.Name
		}
	}
	return ""
}

// findHatName returns name from hat found by id in provided list
func (j *Driver) findHatName(id uint8, hat uint8, list []hat) string {
	for _, lHat := range list {
		if int(id) == lHat.ID && int(hat) == lHat.Hat {
			return lHat.Name
		}
	}
	return ""
}
