// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/beaglebone"
)

func main() {
	beagleboneAdaptor := beaglebone.NewAdaptor()
	button := gpio.NewButtonDriver(beagleboneAdaptor, "P8_9")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("button pressed")
		})

		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("button released")
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{beagleboneAdaptor},
		[]gobot.Device{button},
		work,
	)

	robot.Start()
}
