package gobot

import (
	"errors"
	"fmt"
	"log"
)

// JSONDevice is a JSON representation of a Gobot Device.
type JSONDevice struct {
	Name       string   `json:"name"`
	Driver     string   `json:"driver"`
	Connection string   `json:"connection"`
	Commands   []string `json:"commands"`
}

type Device DriverInterface

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
func (d *devices) Start() error {
	var err error
	log.Println("Starting devices...")
	for _, device := range *d {
		info := "Starting device " + device.Name()
		if device.Pin() != "" {
			info = info + " on pin " + device.Pin()
		}
		log.Println(info + "...")
		err = device.Start()
		if err != nil {
			err = errors.New(fmt.Sprintf("Could not start device: %v", err))
			break
		}
	}
	return err
}

// Halt stop all the devices.
func (d *devices) Halt() {
	for _, device := range *d {
		device.Halt()
	}
}
