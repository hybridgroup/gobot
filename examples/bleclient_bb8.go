//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass the Bluetooth address or name as the first param:

	go run examples/bb8.go BB-1234

 NOTE: sudo is required to use BLE in Linux
*/

//nolint:gosec // ok here
package main

import (
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/sphero"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])
	bb8 := sphero.NewBB8Driver(bleAdaptor)

	work := func() {
		gobot.Every(1*time.Second, func() {
			r := uint8(gobot.Rand(255))
			g := uint8(gobot.Rand(255))
			b := uint8(gobot.Rand(255))
			bb8.SetRGB(r, g, b)
		})
	}

	robot := gobot.NewRobot("bbBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{bb8},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
