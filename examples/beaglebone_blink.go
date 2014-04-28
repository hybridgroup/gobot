package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/beaglebone"
	"github.com/hybridgroup/gobot/gpio"
)

func main() {
	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor()
	beagleboneAdaptor.Name = "beaglebone"

	led := gpio.NewLedDriver(beagleboneAdaptor)
	led.Name = "led"
	led.Pin = "P9_12"

	work := func() {
		gobot.Every("1s", func() {
			led.Toggle()
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{beagleboneAdaptor},
		Devices:     []gobot.Device{led},
		Work:        work,
	}

	robot.Start()
}
