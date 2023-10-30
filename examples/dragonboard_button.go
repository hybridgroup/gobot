//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/dragonboard"
)

func main() {
	dragonAdaptor := dragonboard.NewAdaptor()
	button := gpio.NewButtonDriver(dragonAdaptor, "GPIO_A")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("button pressed")
		})

		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("button released")
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{dragonAdaptor},
		[]gobot.Device{button},
		work,
	)

	robot.Start()
}
