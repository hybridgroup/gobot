// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

// Example of a simple led toggle without the initialization of
// the entire gobot framework.
// This might be useful if you want to use gobot as another
// golang library to interact with sensors and other devices.
func main() {
	f := firmata.NewAdaptor("/dev/ttyACM0")
	f.Connect()

	led := gpio.NewLedDriver(f, "13")
	led.Start()

	for {
		led.Toggle()
		time.Sleep(1000 * time.Millisecond)
	}
}
