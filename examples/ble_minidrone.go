package main

import (
	"os"
	"fmt"
	"time"

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
			fmt.Printf("battery: %d\n", data)
		})

		gobot.On(drone.Event("status"), func(data interface{}) {
			fmt.Printf("status: %d\n", data)
		})

		gobot.On(drone.Event("flying"), func(data interface{}) {
			fmt.Println("flying!")
			gobot.After(5*time.Second, func() {
				fmt.Println("landing...")
				drone.Land()
			})
		})

		gobot.On(drone.Event("landed"), func(data interface{}) {
			fmt.Println("landed.")
		})

		drone.TakeOff()
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{drone},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
