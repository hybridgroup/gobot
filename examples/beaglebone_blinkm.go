package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/beaglebone"
	"github.com/hybridgroup/gobot/i2c"
)

func main() {
	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor()
	beagleboneAdaptor.Name = "beaglebone"

	blinkm := i2c.NewBlinkMDriver(beagleboneAdaptor)
	blinkm.Name = "blinkm"

	work := func() {
		gobot.Every("3s", func() {
			r := byte(gobot.Rand(255))
			g := byte(gobot.Rand(255))
			b := byte(gobot.Rand(255))
			blinkm.Rgb(r, g, b)
			fmt.Println("color", blinkm.Color())
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{beagleboneAdaptor},
		Devices:     []gobot.Device{blinkm},
		Work:        work,
	}

	robot.Start()
}
