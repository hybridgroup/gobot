//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/beagleboard/beaglebone"
)

func main() {
	beagleboneAdaptor := beaglebone.NewAdaptor()
	button := gpio.NewButtonDriver(beagleboneAdaptor, "P8_09")

	work := func() {
		_ = button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("button pressed")
		})

		_ = button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("button released")
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{beagleboneAdaptor},
		[]gobot.Device{button},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
