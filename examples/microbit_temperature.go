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
	ubit := microbit.NewTemperatureDriver(bleAdaptor)

	work := func() {
		ubit.On(microbit.Temperature, func(data interface{}) {
			fmt.Println("Temperature", data)
		})
	}

	robot := gobot.NewRobot("thermoBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ubit},
		work,
	)

	robot.Start()
}
