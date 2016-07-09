package main

import (
	"os"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ble"
)

func main() {
	gbot := gobot.NewGobot()

	bleAdaptor := ble.NewBLEClientAdaptor("ble", os.Args[1])
	ollie := ble.NewSpheroOllieDriver(bleAdaptor, "ollie")

	work := func() {
		ollie.SetRGB(0, 255, 0)
	}

	robot := gobot.NewRobot("ollieBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ollie},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
