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
	"fmt"
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
		if err := drone.TakeOff(); err != nil {
			fmt.Println(err)
		}

		gobot.After(5*time.Second, func() {
			if err := drone.Land(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("minidrone",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{drone},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
