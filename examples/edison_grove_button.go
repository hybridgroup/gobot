//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/intel-iot/edison"
)

func main() {
	e := edison.NewAdaptor()
	button := gpio.NewButtonDriver(e, "2")

	work := func() {
		_ = button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("On!")
		})

		_ = button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("Off!")
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{e},
		[]gobot.Device{button},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
