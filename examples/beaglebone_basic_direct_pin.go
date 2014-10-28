package main

import (
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {

	// Use Gobot to control BeagleBone's digital pins directly

	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
	gpioPin := gpio.NewDirectPinDriver(beagleboneAdaptor, "myDevice", "P9_12")

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
