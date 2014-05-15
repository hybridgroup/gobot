package gobot

import (
	"errors"
	"log"
	"reflect"
	"time"
)

type Device interface {
	Start() bool
	Halt() bool
}

type JsonDevice struct {
	Name       string          `json:"name"`
	Driver     string          `json:"driver"`
	Connection *JsonConnection `json:"connection"`
	Commands   []string        `json:"commands"`
}

type device struct {
	Name     string          `json:"-"`
	Type     string          `json:"-"`
	Interval time.Duration   `json:"-"`
	Robot    *Robot          `json:"-"`
	Driver   DriverInterface `json:"-"`
}

type devices []*device

// Start() starts all the devices.
func (d devices) Start() error {
	var err error
	log.Println("Starting devices...")
	for _, device := range d {
		log.Println("Starting device " + device.Name + "...")
		if device.Start() == false {
			err = errors.New("Could not start device")
			break
		}
	}
	return err
}

// Halt() stop all the devices.
func (d devices) Halt() {
	for _, device := range d {
		device.Halt()
	}
}

func NewDevice(driver DriverInterface, r *Robot) *device {
	d := new(device)
	s := reflect.ValueOf(driver).Type().String()
	d.Type = s[1:len(s)]
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

func (d *device) Halt() bool {
	log.Println("Device " + d.Name + " halted")
	return d.Driver.Halt()
}

func (d *device) Commands() interface{} {
	return FieldByNamePtr(d.Driver, "Commands").Interface()
}

func (d *device) ToJson() *JsonDevice {
	jsonDevice := new(JsonDevice)
	jsonDevice.Name = d.Name
	jsonDevice.Driver = d.Type
	jsonDevice.Connection = d.Robot.Connection(FieldByNamePtr(FieldByNamePtr(d.Driver, "Adaptor").
		Interface().(AdaptorInterface), "Name").
		Interface().(string)).ToJson()
	jsonDevice.Commands = FieldByNamePtr(d.Driver, "Commands").Interface().([]string)
	return jsonDevice
}
