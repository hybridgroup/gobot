// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/microbit"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	ubit := microbit.NewMagnetometerDriver(bleAdaptor)

	work := func() {
		ubit.On(microbit.Magnetometer, func(data interface{}) {
			fmt.Println("Magnetometer", data)
		})
	}

	robot := gobot.NewRobot("magnetoBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ubit},
		work,
	)

	robot.Start()
}
