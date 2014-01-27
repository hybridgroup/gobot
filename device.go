package gobot

import (
	"log"
)

type device struct {
	Name     string          `json:"name"`
	Interval string          `json:"-"`
	Robot    *Robot          `json:"-"`
	Driver   DriverInterface `json:"driver"`
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
	log.Println("Device " + d.Name + " started")
	return d.Driver.Start()
}

func (d *device) Commands() interface{} {
	return FieldByNamePtr(d.Driver, "Commands").Interface()
}
