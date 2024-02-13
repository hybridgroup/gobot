//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass the Bluetooth address or name as the first param:

	go run examples/ble_multiple_info.go BB-1234 BB-1235

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
	info := ble.NewDeviceInformationDriver(bleAdaptor)

	work := func() {
		modelNo, err := info.GetModelNumber()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Model number:", modelNo)

		fwRev, err := info.GetFirmwareRevision()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Firmware rev:", fwRev)

		hwRev, err := info.GetHardwareRevision()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Hardware rev:", hwRev)

		manuName, err := info.GetManufacturerName()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Manufacturer name:", manuName)

		pid, err := info.GetPnPId()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("PnPId:", pid)
	}

	robot := gobot.NewRobot("bot "+port,
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{info},
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
