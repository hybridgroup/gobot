package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
	"time"
)

func main() {
	master := gobot.NewGobot()

	spheros := []string{
		"/dev/rfcomm0",
		"/dev/rfcomm1",
		"/dev/rfcomm2",
		"/dev/rfcomm3",
	}

	for s := range spheros {
		spheroAdaptor := sphero.NewSpheroAdaptor("Sphero", spheros[s])

		spheroDriver := sphero.NewSpheroDriver(spheroAdaptor, "Sphero"+spheros[s])

		work := func() {
			spheroDriver.Stop()

			gobot.On(spheroDriver.Events["Collision"], func(data interface{}) {
				fmt.Println("Collision Detected!")
			})

			gobot.Every(1*time.Second, func() {
				spheroDriver.Roll(100, uint16(gobot.Rand(360)))
			})
			gobot.Every(3*time.Second, func() {
				spheroDriver.SetRGB(uint8(gobot.Rand(255)), uint8(gobot.Rand(255)), uint8(gobot.Rand(255)))
			})
		}

		master.Robots = append(master.Robots,
			gobot.NewRobot("sphero", []gobot.Connection{spheroAdaptor}, []gobot.Device{spheroDriver}, work))
	}

	master.Start()
}
