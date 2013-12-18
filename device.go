package gobot

import (
	"fmt"
)

type device struct {
	Name     string
	Interval string `json:"-"`
	Robot    *Robot `json:"-"`
	Driver   DriverInterface
	Params   map[string]string `json:"-"`
}

type Device interface {
	Start() bool
}

func NewDevice(driver DriverInterface, r *Robot) *device {
	d := new(device)
	d.Name = FieldByNamePtr(driver, "Name").String()
	d.Robot = r
	if FieldByNamePtr(driver, "Interval").String() == "" {
		FieldByNamePtr(driver, "Interval").SetString("0.1s")
	}
	d.Driver = driver
	return d
}

func (d *device) Start() bool {
	fmt.Println("Device " + d.Name + " started")
	d.Driver.Start()
	return true
}

func (d *device) Commands() interface{} {
	return FieldByNamePtr(d.Driver, "Commands").Interface()
}
