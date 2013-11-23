package gobot

import "fmt"

type Driver struct {
	Interval string
	Pin      string
	Name     string
	Params   map[string]string
	Events   map[string]chan interface{} `json:"-"`
}

func NewDriver(d Driver) Driver {
	return d
}

func (d *Driver) Start() {
	fmt.Println("Starting driver " + d.Name + "...")
}
