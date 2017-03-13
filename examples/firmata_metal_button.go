// +build example
//
// Do not build by default.

// TO RUN:
//	firmata_metal_button <PORT>
//
// EXAMPLE:
//	go run ./examples/firmata_metal_button /dev/ttyACM0
//
package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	f := firmata.NewAdaptor(os.Args[1])
	f.Connect()

	led := gpio.NewLedDriver(f, "2")
	led.Start()
	led.Off()

	button := gpio.NewButtonDriver(f, "3")
	button.Start()

	buttonEvents := button.Subscribe()
	for event := range buttonEvents {
		fmt.Println("Event:", event.Name, event.Data)
		if event.Name == gpio.ButtonPush {
			led.Toggle()
		}
	}
}
