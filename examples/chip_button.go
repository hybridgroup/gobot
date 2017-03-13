// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/chip"
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
