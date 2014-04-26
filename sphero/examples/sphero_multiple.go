package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-sphero"
)

func main() {
	master := gobot.GobotMaster()

	spheros := []string{
		"/dev/rfcomm0",
		"/dev/rfcomm1",
		"/dev/rfcomm2",
		"/dev/rfcomm3",
	}

	for s := range spheros {
		spheroAdaptor := new(gobotSphero.SpheroAdaptor)
		spheroAdaptor.Name = "Sphero"
		spheroAdaptor.Port = spheros[s]

		sphero := gobotSphero.NewSphero(spheroAdaptor)
		sphero.Name = "Sphero" + spheros[s]
		sphero.Interval = "0.5s"

		work := func() {
			sphero.Stop()

			gobot.On(sphero.Events["Collision"], func(data interface{}) {
				fmt.Println("Collision Detected!")
			})

			gobot.Every("1s", func() {
				sphero.Roll(100, uint16(gobot.Rand(360)))
			})
			gobot.Every("3s", func() {
				sphero.SetRGB(uint8(gobot.Rand(255)), uint8(gobot.Rand(255)), uint8(gobot.Rand(255)))
			})
		}

		master.Robots = append(master.Robots, &gobot.Robot{
			Connections: []gobot.Connection{spheroAdaptor},
			Devices:     []gobot.Device{sphero},
			Work:        work,
		})
	}

	master.Start()
}
