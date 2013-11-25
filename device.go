package gobot

import (
	"fmt"
	"reflect"
)

type Device struct {
	Name     string
	Interval string
	Robot    *Robot `json:"-"`
	Driver   interface{}
	Params   map[string]string
}

func NewDevice(driver interface{}, r *Robot) *Device {
	d := new(Device)
	d.Name = reflect.ValueOf(driver).Elem().FieldByName("Name").String()
	d.Robot = r
	d.Driver = driver
	return d
}

func (d *Device) Start() {
	fmt.Println("Device " + d.Name + " started")
	r := reflect.ValueOf(d.Driver).MethodByName("StartDriver")
	if r.IsValid() {
		r.Call([]reflect.Value{})
	}
}

func (d *Device) Commands() interface{} {
	return reflect.ValueOf(d.Driver).Elem().FieldByName("Commands").Interface()
}
