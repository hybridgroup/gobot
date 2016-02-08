package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/chip"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	chipAdaptor := chip.NewChipAdaptor("chip")
	button := gpio.NewButtonDriver(chipAdaptor, "button", "XIO-P0")

	work := func() {
		gobot.On(button.Event("push"), func(data interface{}) {
			fmt.Println("button pressed")
		})

		gobot.On(button.Event("release"), func(data interface{}) {
			fmt.Println("button released")
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{chipAdaptor},
		[]gobot.Device{button},
		work,
	)
	gbot.AddRobot(robot)
	gbot.Start()
}
