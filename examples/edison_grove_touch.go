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
	touch := gpio.NewButtonDriver(e, "2")

	work := func() {
		_ = touch.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("On!")
		})

		_ = touch.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("Off!")
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{e},
		[]gobot.Device{touch},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
