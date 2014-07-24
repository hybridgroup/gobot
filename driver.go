package gobot

import (
	"fmt"
	"time"
)

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

func (d *Driver) Adaptor() AdaptorInterface {
	return d.adaptor
}

func (d *Driver) SetInterval(t time.Duration) {
	d.interval = t
}

func (d *Driver) Interval() time.Duration {
	return d.interval
}

func (d *Driver) SetName(s string) {
	d.name = s
}

func (d *Driver) Name() string {
	return d.name
}

func (d *Driver) Pin() string {
	return d.pin
}

func (d *Driver) SetPin(pin string) {
	d.pin = pin
}

func (d *Driver) Type() string {
	return d.driverType
}

func (d *Driver) Events() map[string]*Event {
	return d.events
}

func (d *Driver) Event(name string) *Event {
	e, ok := d.events[name]
	if ok {
		return e
	} else {
		panic(fmt.Sprintf("Unknown Driver Event: %v", name))
	}
}

func (d *Driver) AddEvent(name string) {
	d.events[name] = NewEvent()
}

func (d *Driver) Command(name string) func(map[string]interface{}) interface{} {
	return d.commands[name]
}

func (d *Driver) Commands() map[string]func(map[string]interface{}) interface{} {
	return d.commands
}

func (d *Driver) AddCommand(name string, f func(map[string]interface{}) interface{}) {
	d.commands[name] = f
}

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
