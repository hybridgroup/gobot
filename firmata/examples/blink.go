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
	led.Pin = "13"

	work := func() {
		gobot.Every("1s", func() {
			led.Toggle()
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmata},
		Devices:     []gobot.Device{led},
		Work:        work,
	}

	robot.Start()
}
