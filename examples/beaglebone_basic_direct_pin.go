//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/beagleboard/beaglebone"
)

func main() {
	// Use Gobot to control BeagleBone's digital pins directly
	beagleboneAdaptor := beaglebone.NewAdaptor()
	gpioPin := gpio.NewDirectPinDriver(beagleboneAdaptor, "P9_12")

	// Initialize the internal representation of the pinout
	if err := beagleboneAdaptor.Connect(); err != nil {
		fmt.Println(err)
	}

	// Cast to byte because we are returning an int from a function
	// and not passing in an int literal.
	if err := gpioPin.DigitalWrite(byte(myStateFunction())); err != nil {
		fmt.Println(err)
	}
}

// myStateFunction determines what the GPIO state should be
func myStateFunction() int {
	return 1
}
