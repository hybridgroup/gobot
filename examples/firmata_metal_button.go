package main

import (
	"fmt"

	"github.com/hybridgroup/gobot/drivers/gpio"
	"github.com/hybridgroup/gobot/platforms/firmata"
)

func main() {
	f := firmata.NewAdaptor("/dev/ttyACM0")
	f.Connect()

	led := gpio.NewLedDriver(f, "13")
	led.Start()
	led.Off()

	button := gpio.NewButtonDriver(f, "5")
	button.Start()

	buttonEvents := button.Subscribe()
	for {
		select {
		case event := <-buttonEvents:
			fmt.Println("Event:", event.Name, event.Data)
			if event.Name == gpio.ButtonPush {
				led.Toggle()
			}
		}
	}
}
