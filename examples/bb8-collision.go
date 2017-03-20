// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/sphero/bb8"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	bb := bb8.NewDriver(bleAdaptor)

	work := func() {

		bb.On("collision", func(data interface{}) {
			fmt.Printf("collision detected = %+v \n", data)
			bb.SetRGB(255, 0, 0)
		})

		bb.SetRGB(0, 255, 0)
		bb.Roll(80, 0)
	}

	robot := gobot.NewRobot("bb8",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{bb},
		work,
	)

	robot.Start()

}
