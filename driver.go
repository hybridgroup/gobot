package gobot

import "fmt"

type Driver struct {
  Interval string
  Pin string
  Name string
  Params map[string]string
  Events map[string]chan interface{}
}

func NewDriver(d Driver) Driver {
  return d
}

func (d *Driver) Connection() *interface{}{
  return new(interface{})
}

func (d *Driver) Start() {
  fmt.Println("Starting driver " +  d.Name + "...")
}
