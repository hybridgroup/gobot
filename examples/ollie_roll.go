//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/ble"
	"gobot.io/x/gobot/v2/platforms/sphero/ollie"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	ollie := ollie.NewDriver(bleAdaptor)

	work := func() {
		ollie.SetRGB(255, 0, 255)
		gobot.Every(3*time.Second, func() {
			ollie.Roll(40, uint16(gobot.Rand(360)))
		})
	}

	robot := gobot.NewRobot("ollieBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ollie},
		work,
	)

	robot.Start()
}
