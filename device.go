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
	setInterval(time.Duration)
	getInterval() time.Duration
	setName(string)
	getName() string
}

type JSONDevice struct {
	Name       string          `json:"name"`
	Driver     string          `json:"driver"`
	Connection *JSONConnection `json:"connection"`
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
	d.Name = driver.getName()
	d.Robot = r
	if driver.getInterval() == 0 {
		driver.setInterval(10 * time.Millisecond)
	}
	d.Driver = driver
	return d
}

func (d *device) setInterval(t time.Duration) {
	d.Interval = t
}

func (d *device) getInterval() time.Duration {
	return d.Interval
}

func (d *device) setName(s string) {
	d.Name = s
}

func (d *device) getName() string {
	return d.Name
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

func (d *device) ToJSON() *JSONDevice {
	return &JSONDevice{
		Name:   d.Name,
		Driver: d.Type,
		Connection: d.Robot.Connection(FieldByNamePtr(FieldByNamePtr(d.Driver, "Adaptor").
			Interface().(AdaptorInterface), "Name").
			Interface().(string)).ToJSON(),
		Commands: FieldByNamePtr(d.Driver, "Commands").Interface().([]string),
	}
}
