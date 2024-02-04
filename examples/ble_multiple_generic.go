//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass the Bluetooth address or name as the first param:

	go run examples/ble_multiple_generic.go BB-1234 BB-1235

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func NewSwarmBot(port string) *gobot.Robot {
	bleAdaptor := bleclient.NewAdaptor(port)
	access := ble.NewGenericAccessDriver(bleAdaptor)

	work := func() {
		fmt.Println("Device name:", access.GetDeviceName())
		fmt.Println("Appearance:", access.GetAppearance())
	}

	robot := gobot.NewRobot("bot "+port,
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{access},
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
