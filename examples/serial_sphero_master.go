//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/serial"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
	master := gobot.NewMaster()

	spheros := map[string]string{
		"Sphero-BPO": "/dev/rfcomm0",
	}

	for name, port := range spheros {
		spheroAdaptor := serialport.NewAdaptor(port)
		spheroDriver := serial.NewSpheroDriver(spheroAdaptor)

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
				sphero := master.Robot("Sphero-BPO").Device("sphero").(*serial.SpheroDriver)
				sphero.SetRGB(uint8(gobot.Rand(255)),
					uint8(gobot.Rand(255)),
					uint8(gobot.Rand(255)),
				)
			})
		},
	)

	master.AddRobot(robot)

	if err := master.Start(); err != nil {
		panic(err)
	}
}
