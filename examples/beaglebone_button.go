package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
	button := gpio.NewButtonDriver(beagleboneAdaptor, "button", "P8_9")

	work := func() {
		gobot.On(button.Events["push"], func(data interface{}) {
			fmt.Println("button pressed")
		})

		gobot.On(button.Events["release"], func(data interface{}) {
			fmt.Println("button released")
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("buttonBot", []gobot.Connection{beagleboneAdaptor}, []gobot.Device{button}, work))
	gbot.Start()
}
