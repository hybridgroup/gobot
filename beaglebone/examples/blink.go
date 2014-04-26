package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-beaglebone"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {
	beaglebone := new(gobotBeaglebone.Beaglebone)
	beaglebone.Name = "beaglebone"

	led := gobotGPIO.NewLed(beaglebone)
	led.Name = "led"
	led.Pin = "P9_12"

	work := func() {
		gobot.Every("1s", func() { led.Toggle() })
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{beaglebone},
		Devices:     []gobot.Device{led},
		Work:        work,
	}

	robot.Start()
}
