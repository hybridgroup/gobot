package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-sphero"
)

func main() {
	master := gobot.GobotMaster()
	gobot.Api(master)

	spheros := map[string]string{
		"Sphero-BPO": "127.0.0.1:4560",
	}

	for name, port := range spheros {
		spheroAdaptor := new(gobotSphero.SpheroAdaptor)
		spheroAdaptor.Name = "sphero"
		spheroAdaptor.Port = port

		sphero := gobotSphero.NewSphero(spheroAdaptor)
		sphero.Name = "sphero"
		sphero.Interval = "0.5s"

		work := func() {
			sphero.SetRGB(uint8(255), uint8(0), uint8(0))
		}

		master.Robots = append(master.Robots, gobot.Robot{
			Name:        name,
			Connections: []interface{}{spheroAdaptor},
			Devices:     []interface{}{sphero},
			Work:        work,
		})
	}

	master.Start()
}
