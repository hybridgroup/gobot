package gobot

import (
	"fmt"
	"time"
)

// DriverInterface defines Driver expected behaviour
type DriverInterface interface {
	Start() bool
	Halt() bool
	Adaptor() AdaptorInterface
	SetInterval(time.Duration)
	Interval() time.Duration
	SetName(string)
	Name() string
	Pin() string
	SetPin(string)
	Command(string) func(map[string]interface{}) interface{}
	Commands() map[string]func(map[string]interface{}) interface{}
	AddCommand(string, func(map[string]interface{}) interface{})
	Events() map[string]*Event
	Event(string) *Event
	AddEvent(string)
	Type() string
	ToJSON() *JSONDevice
}

type Driver struct {
	adaptor    AdaptorInterface
	interval   time.Duration
	pin        string
	name       string
	commands   map[string]func(map[string]interface{}) interface{}
	events     map[string]*Event
	driverType string
}

// NewDriver returns a new Driver given a name, driverType and optionally accepts:
//
//	string: Pin the driver connects to
//	AdaptorInterface: Adaptor the driver connects to
//	time.Duration: Interval used internally for polling where applicable
//
// driverType is a label used for identification in the api
func NewDriver(name string, driverType string, v ...interface{}) *Driver {
	if name == "" {
		name = fmt.Sprintf("%X", Rand(int(^uint(0)>>1)))
	}

	d := &Driver{
		driverType: driverType,
		name:       name,
		interval:   10 * time.Millisecond,
		commands:   make(map[string]func(map[string]interface{}) interface{}),
		events:     make(map[string]*Event),
		adaptor:    nil,
		pin:        "",
	}

	for i := range v {
		switch v[i].(type) {
		case string:
			d.pin = v[i].(string)
		case AdaptorInterface:
			d.adaptor = v[i].(AdaptorInterface)
		case time.Duration:
			d.interval = v[i].(time.Duration)
		}
	}

	return d
}

// Adaptor returns driver adaptor
func (d *Driver) Adaptor() AdaptorInterface {
	return d.adaptor
}

// SetInterval defines driver interval duration.
func (d *Driver) SetInterval(t time.Duration) {
	d.interval = t
}

// Interval current driver interval duration
func (d *Driver) Interval() time.Duration {
	return d.interval
}

// SetName sets driver name.
func (d *Driver) SetName(s string) {
	d.name = s
}

// Name returns driver name.
func (d *Driver) Name() string {
	return d.name
}

// Pin returns driver pin
func (d *Driver) Pin() string {
	return d.pin
}

// SetPin defines driver pin
func (d *Driver) SetPin(pin string) {
	d.pin = pin
}

// Type returns driver type
func (d *Driver) Type() string {
	return d.driverType
}

// Events returns driver events map
func (d *Driver) Events() map[string]*Event {
	return d.events
}

// Event returns an event by name if exists
func (d *Driver) Event(name string) *Event {
	e, ok := d.events[name]
	if ok {
		return e
	} else {
		panic(fmt.Sprintf("Unknown Driver Event: %v", name))
	}
}

// AddEvents adds a new event by name
func (d *Driver) AddEvent(name string) {
	d.events[name] = NewEvent()
}

// Command retrieves a command by name
func (d *Driver) Command(name string) func(map[string]interface{}) interface{} {
	return d.commands[name]
}

// Commands returns a map of driver commands
func (d *Driver) Commands() map[string]func(map[string]interface{}) interface{} {
	return d.commands
}

// AddCommand links specified command name to `f`
func (d *Driver) AddCommand(name string, f func(map[string]interface{}) interface{}) {
	d.commands[name] = f
}

// ToJSON returns JSON Driver represnentation including adaptor and commands
func (d *Driver) ToJSON() *JSONDevice {
	jsonDevice := &JSONDevice{
		Name:       d.Name(),
		Driver:     d.Type(),
		Commands:   []string{},
		Connection: "",
	}

	if d.Adaptor() != nil {
		jsonDevice.Connection = d.Adaptor().ToJSON().Name
	}

	for command := range d.Commands() {
		jsonDevice.Commands = append(jsonDevice.Commands, command)
	}

	return jsonDevice
}
