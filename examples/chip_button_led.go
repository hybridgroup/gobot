package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/chip"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	chipAdaptor := chip.NewChipAdaptor("chip")
	button := gpio.NewButtonDriver(chipAdaptor, "button", "XIO-P6")
	led := gpio.NewLedDriver(chipAdaptor, "led", "XIO-P7")

	work := func() {
		gobot.On(button.Event("push"), func(data interface{}) {
			led.On()
		})

		gobot.On(button.Event("release"), func(data interface{}) {
			led.Off()
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{chipAdaptor},
		[]gobot.Device{button, led},
		work,
	)
	gbot.AddRobot(robot)
	gbot.Start()
}
