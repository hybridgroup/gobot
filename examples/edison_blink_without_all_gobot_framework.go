package main

import (
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
	"time"
)

// Example of a simple led toggle without the initialization of
// the entire gobot framework.
// This might be useful if you want to use gobot as another
// golang library to interact with sensors and other devices.
func main() {
	e := edison.NewEdisonAdaptor("edison")
	led := gpio.NewLedDriver(e, "led", "13")
	e.Connect()
	led.Start()
	for {
		led.Toggle()
		time.Sleep(1000 * time.Millisecond)
	}
}
