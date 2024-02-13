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
		devName, err := access.GetDeviceName()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Device name:", devName)

		appearance, err := access.GetAppearance()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Appearance:", appearance)
	}

	robot := gobot.NewRobot("bot "+port,
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{access},
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
