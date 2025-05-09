//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/intel-iot/edison"
)

func main() {
	e := edison.NewAdaptor()
	if err := e.Connect(); err != nil {
		fmt.Println(err)
	}

	led := gpio.NewLedDriver(e, "13")
	if err := led.Start(); err != nil {
		fmt.Println(err)
	}
	if err := led.Off(); err != nil {
		fmt.Println(err)
	}

	button := gpio.NewButtonDriver(e, "5")
	if err := button.Start(); err != nil {
		fmt.Println(err)
	}

	buttonEvents := button.Subscribe()
	for event := range buttonEvents {
		fmt.Println("Event:", event.Name, event.Data)
		if event.Name == gpio.ButtonPush {
			if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
		}
	}
}
