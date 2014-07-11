package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

func main() {
	gbot := gobot.NewGobot()

	spheros := map[string]string{
		"Sphero-BPO": "/dev/rfcomm0",
	}

	for name, port := range spheros {
		spheroAdaptor := sphero.NewSpheroAdaptor("sphero", port)

		spheroDriver := sphero.NewSpheroDriver(spheroAdaptor, "sphero")

		work := func() {
			spheroDriver.SetRGB(uint8(255), uint8(0), uint8(0))
		}

		robot := gobot.NewRobot(name,
			[]gobot.Connection{spheroAdaptor},
			[]gobot.Device{spheroDriver},
			work,
		)

		gbot.AddRobot(robot)
	}

	robot := gobot.NewRobot("",
		func() {
			gobot.Every(1*time.Second, func() {
				sphero := gbot.Robot("Sphero-BPO").Device("sphero").(*sphero.SpheroDriver)
				sphero.SetRGB(uint8(gobot.Rand(255)),
					uint8(gobot.Rand(255)),
					uint8(gobot.Rand(255)),
				)
			})
		},
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
