// +build example
//
// Do not build by default.

package main

import (
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/beaglebone"
)

func main() {
	// Use Gobot to control BeagleBone's digital pins directly
	beagleboneAdaptor := beaglebone.NewAdaptor()
	gpioPin := gpio.NewDirectPinDriver(beagleboneAdaptor, "P9_12")

	// Initialize the internal representation of the pinout
	beagleboneAdaptor.Connect()

	// Cast to byte because we are returning an int from a function
	// and not passing in an int literal.
	gpioPin.DigitalWrite(byte(myStateFunction()))
}

// myStateFunction determines what the GPIO state should be
func myStateFunction() int {
	return 1
}
