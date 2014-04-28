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
	led.Pin = "P9_14"

	work := func() {
		brightness := uint8(0)
		fade_amount := uint8(5)

		gobot.Every("0.1s", func() {
			led.Brightness(brightness)
			brightness = brightness + fade_amount
			if brightness == 0 || brightness == 255 {
				fade_amount = -fade_amount
			}
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{beagleboneAdaptor},
		Devices:     []gobot.Device{led},
		Work:        work,
	}

	robot.Start()
}
