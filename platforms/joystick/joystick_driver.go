package joystick

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	js "github.com/0xcafed00d/joystick"
	"gobot.io/x/gobot/v2"
)

const (
	// Dualshock3 joystick configuration.
	Dualshock3 = "dualshock3"

	// Dualshock4 joystick configuration.
	Dualshock4 = "dualshock4"

	// Dualsense joystick configuration.
	Dualsense = "dualsense"

	// TFlightHotasX flight stick configuration.
	TFlightHotasX = "tflightHotasX"

	// Configuration for Xbox 360 controller.
	Xbox360 = "xbox360"

	// Xbox360RockBandDrums controller configuration.
	Xbox360RockBandDrums = "xbox360RockBandDrums"

	// Configuration for the Xbox One controller.
	XboxOne = "xboxOne"

	// Nvidia Shield TV Controller
	Shield = "shield"

	// Nintendo Switch Joycon Controller Pair
	NintendoSwitchPair = "joyconPair"
)

// Driver represents a joystick
type Driver struct {
	name        string
	interval    time.Duration
	connection  gobot.Connection
	configPath  string
	config      joystickConfig
	buttonState map[int]bool
	axisState   map[int]int

	halt chan bool
	gobot.Eventer
}

// pair is a JSON representation of name and id
type pair struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// joystickConfig is a JSON representation of configuration values
type joystickConfig struct {
	Name    string `json:"name"`
	GUID    string `json:"guid"`
	Axis    []pair `json:"axis"`
	Buttons []pair `json:"buttons"`
}

// NewDriver returns a new Driver with a polling interval of
// 10 Milliseconds given a Joystick Adaptor and json button configuration
// file location.
//
// Optionally accepts:
//
//	time.Duration: Interval at which the Driver is polled for new information
func NewDriver(a *Adaptor, config string, v ...time.Duration) *Driver {
	d := &Driver{
		name:        gobot.DefaultName("Joystick"),
		connection:  a,
		Eventer:     gobot.NewEventer(),
		configPath:  config,
		buttonState: make(map[int]bool),
		axisState:   make(map[int]int),

		interval: 10 * time.Millisecond,
		halt:     make(chan bool),
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
//
//	Error error - On button error
//	Events defined in the json button configuration file.
//	They will have the format:
//		[button]_press
//		[button]_release
//		[axis]
func (j *Driver) Start() error {
	if err := j.initConfig(); err != nil {
		return err
	}

	j.initEvents()

	go func() {
		for {
			state, err := j.adaptor().joystick.Read()
			if err != nil {
				j.Publish(j.Event("error"), err)
				break
			}

			// might just be missing a button definition, so keep going
			if err := j.handleButtons(state); err != nil {
				j.Publish(j.Event("error"), err)
			}

			// might just be missing an axis definition, so keep going
			if err := j.handleAxes(state); err != nil {
				j.Publish(j.Event("error"), err)
			}

			select {
			case <-time.After(j.interval):
			case <-j.halt:
				return
			}
		}
	}()

	return nil
}

func (j *Driver) initConfig() error {
	switch j.configPath {
	case Dualshock3:
		j.config = dualshock3Config
	case Dualshock4:
		j.config = dualshock4Config
	case Dualsense:
		j.config = dualsenseConfig
	case TFlightHotasX:
		j.config = tflightHotasXConfig
	case Xbox360:
		j.config = xbox360Config
	case Xbox360RockBandDrums:
		j.config = xbox360RockBandDrumsConfig
	case XboxOne:
		j.config = xboxOneConfig
	case Shield:
		j.config = shieldConfig
	case NintendoSwitchPair:
		j.config = joyconPairConfig
	default:
		err := j.loadFile()
		if err != nil {
			return fmt.Errorf("loadfile error: %w", err)
		}
	}

	return nil
}

func (j *Driver) initEvents() {
	for _, value := range j.config.Buttons {
		j.AddEvent(fmt.Sprintf("%s_press", value.Name))
		j.AddEvent(fmt.Sprintf("%s_release", value.Name))
	}
	for _, value := range j.config.Axis {
		j.AddEvent(value.Name)
	}
}

// Halt stops joystick driver
func (j *Driver) Halt() (err error) {
	j.halt <- true
	return
}

func (j *Driver) handleButtons(state js.State) error {
	for button := 0; button < j.adaptor().joystick.ButtonCount(); button++ {
		buttonPressed := state.Buttons&(1<<uint32(button)) != 0
		if buttonPressed != j.buttonState[button] {
			j.buttonState[button] = buttonPressed
			name := j.findName(uint8(button), j.config.Buttons)
			if name == "" {
				return fmt.Errorf("Unknown button: %v", button)
			}

			if buttonPressed {
				j.Publish(j.Event(fmt.Sprintf("%s_press", name)), nil)
			} else {
				j.Publish(j.Event(fmt.Sprintf("%s_release", name)), nil)
			}
		}
	}

	return nil
}

func (j *Driver) handleAxes(state js.State) error {
	for axis := 0; axis < j.adaptor().joystick.AxisCount(); axis++ {
		name := j.findName(uint8(axis), j.config.Axis)
		if name == "" {
			return fmt.Errorf("Unknown Axis: %v", axis)
		}

		if j.axisState[axis] != state.AxisData[axis] {
			j.axisState[axis] = state.AxisData[axis]
			j.Publish(name, state.AxisData[axis])
		}
	}

	return nil
}

// findName returns name from button or axis found by id in provided list
func (j *Driver) findName(id uint8, list []pair) string {
	for _, value := range list {
		if int(id) == value.ID {
			return value.Name
		}
	}
	return ""
}

// findID returns the ID based on the name from button or axis.
func (j *Driver) findID(name string, list []pair) int {
	for _, value := range list {
		if name == value.Name {
			return value.ID
		}
	}
	return 0
}

// loadFile load the joystick config from a .json file
func (j *Driver) loadFile() error {
	file, e := os.ReadFile(j.configPath)
	if e != nil {
		return e
	}

	var jsontype joystickConfig
	if err := json.Unmarshal(file, &jsontype); err != nil {
		return err
	}

	j.config = jsontype
	return nil
}
