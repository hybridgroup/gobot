//go:build example
// +build example

//
// Do not build by default.

/*
 To run this example, pass the BLE address or BLE name as first param:

 go run examples/ollie_multiple.go 2B-1234 2B-5678

 NOTE: sudo is required to use BLE in Linux
*/

//nolint:gosec // ok here
package main

import (
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/drivers/ble/sphero"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func NewSwarmBot(port string) *gobot.Robot {
	bleAdaptor := bleclient.NewAdaptor(port)
	ollieDriver := sphero.NewOllieDriver(bleAdaptor)

	work := func() {
		gobot.Every(1*time.Second, func() {
			ollieDriver.SetRGB(uint8(gobot.Rand(255)),
				uint8(gobot.Rand(255)),
				uint8(gobot.Rand(255)),
			)
		})
	}

	robot := gobot.NewRobot("ollie "+port,
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ollieDriver},
		work,
	)

	return robot
}

func main() {
	manager := gobot.NewManager()
	api.NewAPI(manager).Start()

	for _, port := range os.Args[1:] {
		bot := NewSwarmBot(port)
		manager.AddRobot(bot)
	}

	if err := manager.Start(); err != nil {
		panic(err)
	}
}
