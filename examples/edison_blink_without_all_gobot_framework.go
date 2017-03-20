// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
)

// Example of a simple led toggle without the initialization of
// the entire gobot framework.
// This might be useful if you want to use gobot as another
// golang library to interact with sensors and other devices.
func main() {
	e := edison.NewAdaptor()
	e.Connect()

	led := gpio.NewLedDriver(e, "13")
	led.Start()

	for {
		led.Toggle()
		time.Sleep(1000 * time.Millisecond)
	}
}
