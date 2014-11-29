package gobot

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

// JSONDevice is a JSON representation of a Gobot Device.
type JSONDevice struct {
	Name       string   `json:"name"`
	Driver     string   `json:"driver"`
	Connection string   `json:"connection"`
	Commands   []string `json:"commands"`
}

func NewJSONDevice(device Device) *JSONDevice {
	jsonDevice := &JSONDevice{
		Name:       device.Name(),
		Driver:     reflect.TypeOf(device).String(),
		Commands:   []string{},
		Connection: "",
	}
	if device.Connection() != nil {
		jsonDevice.Connection = device.Connection().Name()
	}
	for command := range device.(Commander).Commands() {
		jsonDevice.Commands = append(jsonDevice.Commands, command)
	}
	return jsonDevice
}

type Device Driver

type devices []Device

// Len returns devices length
func (d *devices) Len() int {
	return len(*d)
}

// Each calls `f` function each device
func (d *devices) Each(f func(Device)) {
	for _, device := range *d {
		f(device)
	}
}

// Start starts all the devices.
func (d *devices) Start() (errs []error) {
	log.Println("Starting devices...")
	for _, device := range *d {
		info := "Starting device " + device.Name()

		if pinner, ok := device.(Pinner); ok {
			info = info + " on pin " + pinner.Pin()
		}

		log.Println(info + "...")
		if errs = device.Start(); len(errs) > 0 {
			for i, err := range errs {
				errs[i] = errors.New(fmt.Sprintf("Device %q: %v", device.Name(), err))
			}
			return
		}
	}
	return
}

// Halt stop all the devices.
func (d *devices) Halt() (errs []error) {
	for _, device := range *d {
		if derrs := device.Halt(); len(derrs) > 0 {
			for i, err := range derrs {
				derrs[i] = errors.New(fmt.Sprintf("Device %q: %v", device.Name(), err))
			}
			errs = append(errs, derrs...)
		}
	}
	return
}
