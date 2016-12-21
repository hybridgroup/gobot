package gobot

import (
	"log"
	"reflect"

	multierror "github.com/hashicorp/go-multierror"
)

// JSONDevice is a JSON representation of a Device.
type JSONDevice struct {
	Name       string   `json:"name"`
	Driver     string   `json:"driver"`
	Connection string   `json:"connection"`
	Commands   []string `json:"commands"`
}

// NewJSONDevice returns a JSONDevice given a Device.
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
	if commander, ok := device.(Commander); ok {
		for command := range commander.Commands() {
			jsonDevice.Commands = append(jsonDevice.Commands, command)
		}
	}
	return jsonDevice
}

// A Device is an instnace of a Driver
type Device Driver

// Devices represents a collection of Device
type Devices []Device

// Len returns devices length
func (d *Devices) Len() int {
	return len(*d)
}

// Each enumerates through the Devices and calls specified callback function.
func (d *Devices) Each(f func(Device)) {
	for _, device := range *d {
		f(device)
	}
}

// Start calls Start on each Device in d
func (d *Devices) Start() (err error) {
	log.Println("Starting devices...")
	for _, device := range *d {
		info := "Starting device " + device.Name()

		if pinner, ok := device.(Pinner); ok {
			info = info + " on pin " + pinner.Pin()
		}

		log.Println(info + "...")
		if derr := device.Start(); derr != nil {
			err = multierror.Append(err, derr)
		}
	}
	return err
}

// Halt calls Halt on each Device in d
func (d *Devices) Halt() (err error) {
	for _, device := range *d {
		if derr := device.Halt(); derr != nil {
			err = multierror.Append(err, derr)
		}
	}
	return err
}
