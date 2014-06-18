package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
	"time"
)

func main() {
	gbot := gobot.NewGobot()

	adaptor := sphero.NewSpheroAdaptor("Sphero", "/dev/rfcomm0")
	spheroDriver := sphero.NewSpheroDriver(adaptor, "sphero")

	work := func() {
		gobot.On(spheroDriver.Events["Collision"], func(data interface{}) {
			fmt.Println("Collision Detected!")
		})

		gobot.Every(3*time.Second, func() {
			spheroDriver.Roll(30, uint16(gobot.Rand(360)))
		})

		gobot.Every(1*time.Second, func() {
			r := uint8(gobot.Rand(255))
			g := uint8(gobot.Rand(255))
			b := uint8(gobot.Rand(255))
			spheroDriver.SetRGB(r, g, b)
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("sphero", []gobot.Connection{adaptor}, []gobot.Device{spheroDriver}, work))

	gbot.Start()
}
