package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sphero"
)

func main() {
	master := gobot.NewMaster()

	spheros := []string{
		"/dev/rfcomm0",
		"/dev/rfcomm1",
		"/dev/rfcomm2",
		"/dev/rfcomm3",
	}

	for s := range spheros {
		spheroAdaptor := sphero.NewSpheroAdaptor()
		spheroAdaptor.Name = "Sphero"
		spheroAdaptor.Port = spheros[s]

		spheroDriver := sphero.NewSpheroDriver(spheroAdaptor)
		spheroDriver.Name = "Sphero" + spheros[s]
		spheroDriver.Interval = "0.5s"

		work := func() {
			spheroDriver.Stop()

			gobot.On(spheroDriver.Events["Collision"], func(data interface{}) {
				fmt.Println("Collision Detected!")
			})

			gobot.Every("1s", func() {
				spheroDriver.Roll(100, uint16(gobot.Rand(360)))
			})
			gobot.Every("3s", func() {
				spheroDriver.SetRGB(uint8(gobot.Rand(255)), uint8(gobot.Rand(255)), uint8(gobot.Rand(255)))
			})
		}

		master.Robots = append(master.Robots, &gobot.Robot{
			Connections: []gobot.Connection{spheroAdaptor},
			Devices:     []gobot.Device{spheroDriver},
			Work:        work,
		})
	}

	master.Start()
}
