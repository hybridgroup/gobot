// +build example
//
// Do not build by default.

/*
 How to run
 Pass the Bluetooth address or name as the first param:

	go run examples/ble_generic_access.go BB-1234

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	access := ble.NewGenericAccessDriver(bleAdaptor)

	work := func() {
		fmt.Println("Device name:", access.GetDeviceName())
		fmt.Println("Appearance:", access.GetAppearance())
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{access},
		work,
	)

	robot.Start()
}
