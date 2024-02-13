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
	a := api.NewAPI(manager)
	a.Port = "8080"
	a.Start()

	pebbleAdaptor := pebble.NewAdaptor()
	pebbleDriver := pebble.NewDriver(pebbleAdaptor)

	work := func() {
		_ = pebbleDriver.On(pebbleDriver.Event("accel"), func(data interface{}) {
			fmt.Println(data.(string))
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
