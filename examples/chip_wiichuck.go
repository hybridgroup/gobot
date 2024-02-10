//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/chip"
)

func main() {
	chipAdaptor := chip.NewAdaptor()
	wiichuck := i2c.NewWiichuckDriver(chipAdaptor)

	work := func() {
		_ = wiichuck.On(wiichuck.Event("joystick"), func(data interface{}) {
			fmt.Println("joystick", data)
		})

		_ = wiichuck.On(wiichuck.Event("c"), func(data interface{}) {
			fmt.Println("c")
		})

		_ = wiichuck.On(wiichuck.Event("z"), func(data interface{}) {
			fmt.Println("z")
		})
		_ = wiichuck.On(wiichuck.Event("error"), func(data interface{}) {
			fmt.Println("Wiichuck error:", data)
		})
	}

	robot := gobot.NewRobot("chuck",
		[]gobot.Connection{chipAdaptor},
		[]gobot.Device{wiichuck},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
