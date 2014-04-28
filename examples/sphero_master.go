package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sphero"
)

func main() {
	master := gobot.GobotMaster()

	spheros := map[string]string{
		"Sphero-BPO": "/dev/rfcomm0",
	}

	for name, port := range spheros {
		spheroAdaptor := new(sphero.Adaptor)
		spheroAdaptor.Name = "sphero"
		spheroAdaptor.Port = port

		spheroDriver := sphero.NewSpheroDriver(spheroAdaptor)
		spheroDriver.Name = "sphero"
		spheroDriver.Interval = "0.5s"

		work := func() {
			spheroDriver.SetRGB(uint8(255), uint8(0), uint8(0))
		}

		master.Robots = append(master.Robots, &gobot.Robot{
			Name:        name,
			Connections: []gobot.Connection{spheroAdaptor},
			Devices:     []gobot.Device{spheroDriver},
			Work:        work,
		})
	}

	master.Robots = append(master.Robots, &gobot.Robot{
		Work: func() {
			gobot.Every("1s", func() {
				gobot.Call(master.FindRobot("Sphero-BPO").GetDevice("spheroDriver").Driver, "SetRGB", uint8(gobot.Rand(255)), uint8(gobot.Rand(255)), uint8(gobot.Rand(255)))
			})
		},
	})

	master.Start()
}
