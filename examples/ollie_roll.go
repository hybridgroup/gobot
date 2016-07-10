package main

import (
	"os"
	"time"
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ble"
)

func main() {
	gbot := gobot.NewGobot()

	bleAdaptor := ble.NewBLEClientAdaptor("ble", os.Args[1])
	ollie := ble.NewSpheroOllieDriver(bleAdaptor, "ollie")

	work := func() {
		ollie.SetRGB(255, 0, 255)
		gobot.After(1*time.Second, func() {
			fmt.Println("fwd")
			ollie.Roll(60, 0)
		})
		gobot.After(3*time.Second, func() {
			fmt.Println("back")
			ollie.Roll(60, 180)
		})
		gobot.After(5*time.Second, func() {
			fmt.Println("stop")
			ollie.Stop()
		})
	}

	robot := gobot.NewRobot("ollieBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ollie},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
