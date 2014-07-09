package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

func main() {
	gbot := gobot.NewGobot()
	api.NewAPI(gbot).Start()

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
		robot.AddCommand("turn_blue", func(params map[string]interface{}) interface{} {
			spheroDriver.SetRGB(uint8(0), uint8(0), uint8(255))
			return nil
		})

		gbot.AddRobot(robot)
	}

	gbot.Start()
}
