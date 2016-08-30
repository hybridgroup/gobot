package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

func main() {
	gbot := gobot.NewGobot()

	r := raspi.NewRaspiAdaptor("raspi")
	button := gpio.NewButtonDriver(r, "button", "11")
	led := gpio.NewLedDriver(r, "led", "7")

	work := func() {
		button.On(button.Event("push"), func(data interface{}) {
			fmt.Println("button pressed")
			led.On()
		})

		button.On(button.Event("release"), func(data interface{}) {
			fmt.Println("button released")
			led.Off()
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{r},
		[]gobot.Device{button, led},
		work,
	)
	gbot.AddRobot(robot)
	gbot.Start()
}
