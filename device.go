package gobot

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
)

type Device interface {
	Start() bool
	Halt() bool
	setInterval(time.Duration)
	interval() time.Duration
	setName(string)
	name() string
	commands() map[string]func(map[string]interface{}) interface{}
}

type JSONDevice struct {
	Name       string          `json:"name"`
	Driver     string          `json:"driver"`
	Connection *JSONConnection `json:"connection"`
	Commands   []string        `json:"commands"`
}

type device struct {
	Name   string          `json:"-"`
	Type   string          `json:"-"`
	Robot  *Robot          `json:"-"`
	Driver DriverInterface `json:"-"`
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
	if driver.name() == "" {
		driver.setName(fmt.Sprintf("%X", Rand(int(^uint(0)>>1))))
	}
	t := reflect.ValueOf(driver).Type().String()
	if driver.interval() == 0 {
		driver.setInterval(10 * time.Millisecond)
	}
	return &device{
		Type:   t[1:len(t)],
		Name:   driver.name(),
		Robot:  r,
		Driver: driver,
	}
}

func (d *device) setInterval(t time.Duration) {
	d.Driver.setInterval(t)
}

func (d *device) interval() time.Duration {
	return d.Driver.interval()
}

func (d *device) setName(s string) {
	d.Name = s
}

func (d *device) name() string {
	return d.Name
}

func (d *device) commands() map[string]func(map[string]interface{}) interface{} {
	return d.Driver.commands()
}

func (d *device) Commands() map[string]func(map[string]interface{}) interface{} {
	return d.commands()
}

func (d *device) Start() bool {
	log.Println("Device " + d.Name + " started")
	return d.Driver.Start()
}

func (d *device) Halt() bool {
	log.Println("Device " + d.Name + " halted")
	return d.Driver.Halt()
}

func (d *device) ToJSON() *JSONDevice {
	jsonDevice := &JSONDevice{
		Name:   d.Name,
		Driver: d.Type,
		Connection: d.Robot.Connection(FieldByNamePtr(FieldByNamePtr(d.Driver, "Adaptor").
			Interface().(AdaptorInterface), "Name").
			Interface().(string)).ToJSON(),
		Commands: []string{},
	}

	commands := d.commands()
	for command := range commands {
		jsonDevice.Commands = append(jsonDevice.Commands, command)
	}

	return jsonDevice
}
