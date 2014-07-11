package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

func main() {
	gbot := gobot.NewGobot()

	spheros := []string{
		"/dev/rfcomm0",
		"/dev/rfcomm1",
		"/dev/rfcomm2",
		"/dev/rfcomm3",
	}

	for _, port := range spheros {
		spheroAdaptor := sphero.NewSpheroAdaptor("Sphero", port)
		spheroDriver := sphero.NewSpheroDriver(spheroAdaptor, "Sphero"+port)

		work := func() {
			spheroDriver.Stop()

			gobot.On(spheroDriver.Event("collision"), func(data interface{}) {
				fmt.Println("Collision Detected!")
			})

			gobot.Every(1*time.Second, func() {
				spheroDriver.Roll(100, uint16(gobot.Rand(360)))
			})
			gobot.Every(3*time.Second, func() {
				spheroDriver.SetRGB(uint8(gobot.Rand(255)),
					uint8(gobot.Rand(255)),
					uint8(gobot.Rand(255)),
				)
			})
		}

		robot := gobot.NewRobot("sphero",
			[]gobot.Connection{spheroAdaptor},
			[]gobot.Device{spheroDriver},
			work,
		)
		gbot.AddRobot(robot)
	}

	gbot.Start()
}
