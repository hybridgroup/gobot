package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ble"
)

func main() {
	gbot := gobot.NewGobot()

	bleAdaptor := ble.NewBLEAdaptor("ble", "D0:39:72:C9:9E:5A")
	battery := ble.NewBLEBatteryDriver(bleAdaptor, "battery")

	work := func() {
		fmt.Println("Battery level:", battery.GetBatteryLevel())
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{battery},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
