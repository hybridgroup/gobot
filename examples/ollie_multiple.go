//go:build example
// +build example

//
// Do not build by default.

/*
 To run this example, pass the BLE address or BLE name as first param:

 go run examples/ollie_multiple.go 2B-1234 2B-5678

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/platforms/ble"
	"gobot.io/x/gobot/v2/platforms/sphero/ollie"
)

func NewSwarmBot(port string) *gobot.Robot {
	bleAdaptor := ble.NewClientAdaptor(port)
	ollieDriver := ollie.NewDriver(bleAdaptor)

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
	master := gobot.NewMaster()
	api.NewAPI(master).Start()

	for _, port := range os.Args[1:] {
		bot := NewSwarmBot(port)
		master.AddRobot(bot)
	}

	master.Start()
}
