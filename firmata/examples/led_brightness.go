package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-firmata"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {
	firmata := new(gobotFirmata.FirmataAdaptor)
	firmata.Name = "firmata"
	firmata.Port = "/dev/ttyACM0"

	led := gobotGPIO.NewLed(firmata)
	led.Name = "led"
	led.Pin = "3"

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
		Connections: []gobot.Connection{firmata},
		Devices:     []gobot.Device{led},
		Work:        work,
	}

	robot.Start()
}
