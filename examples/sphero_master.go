package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
	"time"
)

func main() {
	master := gobot.NewGobot()

	spheros := map[string]string{
		"Sphero-BPO": "/dev/rfcomm0",
	}

	for name, port := range spheros {
		spheroAdaptor := sphero.NewSpheroAdaptor("sphero", port)

		spheroDriver := sphero.NewSpheroDriver(spheroAdaptor, "sphero")

		work := func() {
			spheroDriver.SetRGB(uint8(255), uint8(0), uint8(0))
		}

		master.Robots = append(master.Robots,
			gobot.NewRobot(name, []gobot.Connection{spheroAdaptor}, []gobot.Device{spheroDriver}, work))
	}

	master.Robots = append(master.Robots, gobot.NewRobot(
		""
		nil,
		nil,
		func() {
			gobot.Every(1*time.Second, func() {
				gobot.Call(master.FindRobot("Sphero-BPO").GetDevice("spheroDriver").Driver, "SetRGB", uint8(gobot.Rand(255)), uint8(gobot.Rand(255)), uint8(gobot.Rand(255)))
			})
		},
	))

	master.Start()
}
