//go:build example
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
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1], bleclient.WithScanTimeout(30*time.Second), bleclient.WithDebug())
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

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{access},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
