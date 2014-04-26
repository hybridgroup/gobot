package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-digispark"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {

	digispark := new(gobotDigispark.DigisparkAdaptor)
	digispark.Name = "Digispark"

	led := gobotGPIO.NewLed(digispark)
	led.Name = "led"
	led.Pin = "0"

	work := func() {
		brightness := uint8(0)
		fade_amount := uint8(15)

		gobot.Every("0.1s", func() {
			led.Brightness(brightness)
			brightness = brightness + fade_amount
			if brightness == 0 || brightness == 255 {
				fade_amount = -fade_amount
			}
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{digispark},
		Devices:     []gobot.Device{led},
		Work:        work,
	}

	robot.Start()
}
