package main

import (
	"fmt"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

func main() {
	e := edison.NewEdisonAdaptor("edison")
	led := gpio.NewLedDriver(e, "led", "13")
	button := gpio.NewButtonDriver(e, "button", "5")

	e.Connect()
	led.Start()
	button.Start()

	led.Off()

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
