package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sphero"
)

func main() {
	adaptor := sphero.NewAdaptor()
	adaptor.Name = "Sphero"
	adaptor.Port = "/dev/rfcomm0"

	spheroDriver := sphero.NewSpheroDriver(adaptor)
	spheroDriver.Name = "sphero"

	work := func() {
		gobot.On(spheroDriver.Events["Collision"], func(data interface{}) {
			fmt.Println("Collision Detected!")
		})

		gobot.Every("3s", func() {
			spheroDriver.Roll(30, uint16(gobot.Rand(360)))
		})

		gobot.Every("1s", func() {
			r := uint8(gobot.Rand(255))
			g := uint8(gobot.Rand(255))
			b := uint8(gobot.Rand(255))
			spheroDriver.SetRGB(r, g, b)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{adaptor},
		Devices:     []gobot.Device{spheroDriver},
		Work:        work,
	}

	robot.Start()
}
