//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/intel-iot/edison"
)

// Example of a simple led toggle without the initialization of
// the entire gobot framework.
// This might be useful if you want to use gobot as another
// golang library to interact with sensors and other devices.
func main() {
	e := edison.NewAdaptor()
	if err := e.Connect(); err != nil {
		fmt.Println(err)
	}

	led := gpio.NewLedDriver(e, "13")
	if err := led.Start(); err != nil {
		fmt.Println(err)
	}

	for {
		if err := led.Toggle(); err != nil {
			fmt.Println(err)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
