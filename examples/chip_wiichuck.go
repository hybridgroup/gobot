// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/chip"
)

func main() {
	chipAdaptor := chip.NewAdaptor()
	wiichuck := i2c.NewWiichuckDriver(chipAdaptor)

	work := func() {
		wiichuck.On(wiichuck.Event("joystick"), func(data interface{}) {
			fmt.Println("joystick", data)
		})

		wiichuck.On(wiichuck.Event("c"), func(data interface{}) {
			fmt.Println("c")
		})

		wiichuck.On(wiichuck.Event("z"), func(data interface{}) {
			fmt.Println("z")
		})
		wiichuck.On(wiichuck.Event("error"), func(data interface{}) {
			fmt.Println("Wiichuck error:", data)
		})
	}

	robot := gobot.NewRobot("chuck",
		[]gobot.Connection{chipAdaptor},
		[]gobot.Device{wiichuck},
		work,
	)

	robot.Start()
}
