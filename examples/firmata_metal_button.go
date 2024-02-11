//go:build example
// +build example

//
// Do not build by default.

// TO RUN:
//
//	firmata_metal_button <PORT>
//
// EXAMPLE:
//
//	go run ./examples/firmata_metal_button /dev/ttyACM0
package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	f := firmata.NewAdaptor(os.Args[1])
	if err := f.Connect(); err != nil {
		fmt.Println(err)
	}

	led := gpio.NewLedDriver(f, "2")
	if err := led.Start(); err != nil {
		fmt.Println(err)
	}
	if err := led.Off(); err != nil {
		fmt.Println(err)
	}

	button := gpio.NewButtonDriver(f, "3")
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
