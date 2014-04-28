package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/beaglebone"
	"github.com/hybridgroup/gobot/gpio"
)

func main() {
	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor()
	beagleboneAdaptor.Name = "beaglebone"

	servo := gpio.NewServoDriver(beagleboneAdaptor)
	servo.Name = "servo"
	servo.Pin = "P9_14"

	work := func() {
		gobot.Every("1s", func() {
			i := uint8(gobot.Rand(180))
			fmt.Println("Turning", i)
			servo.Move(i)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{beagleboneAdaptor},
		Devices:     []gobot.Device{servo},
		Work:        work,
	}

	robot.Start()
}
