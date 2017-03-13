// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/sphero"
)

func main() {
	master := gobot.NewMaster()

	spheros := map[string]string{
		"Sphero-BPO": "/dev/rfcomm0",
	}

	for name, port := range spheros {
		spheroAdaptor := sphero.NewAdaptor(port)
		spheroDriver := sphero.NewSpheroDriver(spheroAdaptor)

		work := func() {
			spheroDriver.SetRGB(uint8(255), uint8(0), uint8(0))
		}

		robot := gobot.NewRobot(name,
			[]gobot.Connection{spheroAdaptor},
			[]gobot.Device{spheroDriver},
			work,
		)

		master.AddRobot(robot)
	}

	robot := gobot.NewRobot("",
		func() {
			gobot.Every(1*time.Second, func() {
				sphero := master.Robot("Sphero-BPO").Device("sphero").(*sphero.SpheroDriver)
				sphero.SetRGB(uint8(gobot.Rand(255)),
					uint8(gobot.Rand(255)),
					uint8(gobot.Rand(255)),
				)
			})
		},
	)

	master.AddRobot(robot)

	master.Start()
}
