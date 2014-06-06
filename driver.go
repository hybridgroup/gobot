package gobot

import "time"

type Driver struct {
	Interval time.Duration               `json:"interval"`
	Pin      string                      `json:"pin"`
	Name     string                      `json:"name"`
	Commands []string                    `json:"commands"`
	Events   map[string]chan interface{} `json:"-"`
}

type DriverInterface interface {
	Start() bool
	Halt() bool
	setInterval(time.Duration)
	getInterval() time.Duration
	setName(string)
	getName() string
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
