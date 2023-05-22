//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/chip"
)

func main() {
	chipAdaptor := chip.NewAdaptor()
	button := gpio.NewButtonDriver(chipAdaptor, "XIO-P0")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("button pressed")
		})

		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("button released")
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{chipAdaptor},
		[]gobot.Device{button},
		work,
	)

	robot.Start()
}
