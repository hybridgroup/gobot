package gobot

import (
	"errors"
	"log"
)

type JSONDevice struct {
	Name       string   `json:"name"`
	Driver     string   `json:"driver"`
	Connection string   `json:"connection"`
	Commands   []string `json:"commands"`
}

type Device DriverInterface

type devices []Device

func (d *devices) Len() int {
	return len(*d)
}

func (d *devices) Each(f func(Device)) {
	for _, device := range *d {
		f(device)
	}
}

// Start() starts all the devices.
func (d *devices) Start() error {
	var err error
	log.Println("Starting devices...")
	for _, device := range *d {
		info := "Starting device " + device.Name()
		if device.Pin() != "" {
			info = info + " on pin " + device.Pin()
		}
		log.Println(info + "...")
		if device.Start() == false {
			err = errors.New("Could not start device")
			break
		}
	}
	return err
}

// Halt() stop all the devices.
func (d *devices) Halt() {
	for _, device := range *d {
		device.Halt()
	}
}
