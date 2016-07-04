package main

import (
	"fmt"
	"time"
	"os"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ble"
)

func main() {
	gbot := gobot.NewGobot()

	bleAdaptor := ble.NewBLEAdaptor("ble", os.Args[1])
	drone := ble.NewBLEMinidroneDriver(bleAdaptor, "drone")

	work := func() {
		battery.Init()
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{drone},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
