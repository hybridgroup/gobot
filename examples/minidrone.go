//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass the Bluetooth name or address as first param:

	go run examples/minidrone.go "Mambo_1234"

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/parrot"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])
	drone := parrot.NewMinidroneDriver(bleAdaptor)

	work := func() {
		drone.TakeOff()

		gobot.After(5*time.Second, func() {
			drone.Land()
		})
	}

	robot := gobot.NewRobot("minidrone",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
