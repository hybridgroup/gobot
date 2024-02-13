//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/platforms/pebble"
)

func main() {
	manager := gobot.NewManager()
	api := api.NewAPI(manager)
	api.Port = "8080"
	api.Start()

	pebbleAdaptor := pebble.NewAdaptor()
	pebbleDriver := pebble.NewDriver(pebbleAdaptor)

	work := func() {
		pebbleDriver.SendNotification("Hello Pebble!")
		_ = pebbleDriver.On(pebbleDriver.Event("button"), func(data interface{}) {
			fmt.Println("Button pushed: " + data.(string))
		})

		_ = pebbleDriver.On(pebbleDriver.Event("tap"), func(data interface{}) {
			fmt.Println("Tap event detected")
		})
	}

	robot := gobot.NewRobot("pebble",
		[]gobot.Connection{pebbleAdaptor},
		[]gobot.Device{pebbleDriver},
		work,
	)

	manager.AddRobot(robot)

	if err := manager.Start(); err != nil {
		panic(err)
	}
}
