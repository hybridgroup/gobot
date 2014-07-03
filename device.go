package gobot

import (
	"errors"
	"log"
)

type JSONDevice struct {
	Name       string          `json:"name"`
	Driver     string          `json:"driver"`
	Connection *JSONConnection `json:"connection"`
	Commands   []string        `json:"commands"`
}

type Device DriverInterface

type devices struct {
	devices []Device
}

func (d *devices) Len() int {
	return len(d.devices)
}

func (d *devices) Add(dev Device) Device {
	d.devices = append(d.devices, dev)
	return dev
}

func (d *devices) Each(f func(Device)) {
	for _, device := range d.devices {
		f(device)
	}
}

// Start() starts all the devices.
func (d devices) Start() error {
	var err error
	log.Println("Starting devices...")
	for _, device := range d.devices {
		log.Println("Starting device " + device.name() + "...")
		if device.Start() == false {
			err = errors.New("Could not start device")
			break
		}
	}
	return err
}

// Halt() stop all the devices.
func (d devices) Halt() {
	for _, device := range d.devices {
		device.Halt()
	}
}
