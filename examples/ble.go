package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ble"
)

func main() {
	gbot := gobot.NewGobot()

	bleAdaptor := ble.NewBLEAdaptor("ble", "D7:99:5A:26:EC:38")
	battery := ble.NewBLEBatteryDriver(bleAdaptor, "battery")

	work := func() {
		fmt.Println("Working...")
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{battery},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
