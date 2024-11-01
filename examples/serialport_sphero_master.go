//go:build example
// +build example

//
// Do not build by default.

//nolint:gosec // ok here
package main

import (
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/serial/sphero"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
	manager := gobot.NewManager()

	spheros := map[string]string{
		"Sphero-BPO": "/dev/rfcomm0",
	}

	for name, port := range spheros {
		spheroAdaptor := serialport.NewAdaptor(port)
		spheroDriver := sphero.NewSpheroDriver(spheroAdaptor)

		work := func() {
			spheroDriver.SetRGB(uint8(255), uint8(0), uint8(0))
		}

		robot := gobot.NewRobot(name,
			[]gobot.Connection{spheroAdaptor},
			[]gobot.Device{spheroDriver},
			work,
		)

		manager.AddRobot(robot)
	}

	robot := gobot.NewRobot("",
		func() {
			gobot.Every(1*time.Second, func() {
				sphero := manager.Robot("Sphero-BPO").Device("sphero").(*sphero.SpheroDriver)
				sphero.SetRGB(uint8(gobot.Rand(255)),
					uint8(gobot.Rand(255)),
					uint8(gobot.Rand(255)),
				)
			})
		},
	)

	manager.AddRobot(robot)

	if err := manager.Start(); err != nil {
		panic(err)
	}
}
