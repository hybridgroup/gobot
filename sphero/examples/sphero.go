package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-sphero"
)

func main() {
	spheroAdaptor := new(gobotSphero.SpheroAdaptor)
	spheroAdaptor.Name = "Sphero"
	spheroAdaptor.Port = "/dev/rfcomm0"

	sphero := gobotSphero.NewSphero(spheroAdaptor)
	sphero.Name = "sphero"

	work := func() {
		gobot.On(sphero.Events["Collision"], func(data interface{}) {
			fmt.Println("Collision Detected!")
		})

		gobot.Every("3s", func() {
			sphero.Roll(30, uint16(gobot.Rand(360)))
		})

		gobot.Every("1s", func() {
			r := uint8(gobot.Rand(255))
			g := uint8(gobot.Rand(255))
			b := uint8(gobot.Rand(255))
			sphero.SetRGB(r, g, b)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{spheroAdaptor},
		Devices:     []gobot.Device{sphero},
		Work:        work,
	}

	robot.Start()
}
