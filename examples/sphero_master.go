package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-sphero"
)

func main() {
	bot := gobot.GobotMaster()

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

		bot.Robots = append(bot.Robots, gobot.Robot{
			Name:        name,
			Connections: []interface{}{spheroAdaptor},
			Devices:     []interface{}{sphero},
			Work:        work,
		})
	}

	bot.Robots = append(bot.Robots, gobot.Robot{
		Work: func() {
			sphero := bot.FindRobot("Sphero-BPO")
			gobot.Every("1s", func() {
				gobot.Call(sphero.GetDevice("sphero").Driver, "SetRGB", uint8(gobot.Rand(255)), uint8(gobot.Rand(255)), uint8(gobot.Rand(255)))
			})
		},
	})

	bot.Start()
}
