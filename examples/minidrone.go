// +build example
//
// Do not build by default.

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/parrot/minidrone"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	drone := minidrone.NewDriver(bleAdaptor)

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
