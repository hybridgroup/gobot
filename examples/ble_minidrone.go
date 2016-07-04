package main

import (
	"os"
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ble"
)

func main() {
	gbot := gobot.NewGobot()

	bleAdaptor := ble.NewBLEAdaptor("ble", os.Args[1])
	drone := ble.NewBLEMinidroneDriver(bleAdaptor, "drone")

	work := func() {
		drone.Init()

		gobot.On(drone.Event("battery"), func(data interface{}) {
			fmt.Printf("battery: %d", data)
		})

		gobot.On(drone.Event("status"), func(data interface{}) {
			fmt.Printf("status: %d", data)
		})
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{drone},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
