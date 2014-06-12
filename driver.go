package gobot

import "time"

type Driver struct {
	Interval time.Duration                                       `json:"interval"`
	Pin      string                                              `json:"pin"`
	Name     string                                              `json:"name"`
	Commands map[string]func(map[string]interface{}) interface{} `json:"commands"`
	Events   map[string]*Event                                   `json:"-"`
}

type DriverInterface interface {
	Start() bool
	Halt() bool
	setInterval(time.Duration)
	getInterval() time.Duration
	setName(string)
	getName() string
	getCommands() map[string]func(map[string]interface{}) interface{}
}

func (d *Driver) setInterval(t time.Duration) {
	d.Interval = t
}

func (d *Driver) getInterval() time.Duration {
	return d.Interval
}

func (d *Driver) setName(s string) {
	d.Name = s
}

func (d *Driver) getName() string {
	return d.Name
}

func (d *Driver) AddCommand(name string, f func(map[string]interface{}) interface{}) {
	d.Commands[name] = f
}

func (d *Driver) getCommands() map[string]func(map[string]interface{}) interface{} {
	return d.Commands
}
