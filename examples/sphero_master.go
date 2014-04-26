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

		sphero := sphero.NewSphero(spheroAdaptor)
		sphero.Name = "sphero"
		sphero.Interval = "0.5s"

		work := func() {
			sphero.SetRGB(uint8(255), uint8(0), uint8(0))
		}

		master.Robots = append(master.Robots, &gobot.Robot{
			Name:        name,
			Connections: []gobot.Connection{spheroAdaptor},
			Devices:     []gobot.Device{sphero},
			Work:        work,
		})
	}

	master.Robots = append(master.Robots, &gobot.Robot{
		Work: func() {
			gobot.Every("1s", func() {
				gobot.Call(master.FindRobot("Sphero-BPO").GetDevice("sphero").Driver, "SetRGB", uint8(gobot.Rand(255)), uint8(gobot.Rand(255)), uint8(gobot.Rand(255)))
			})
		},
	})

	master.Start()
}
