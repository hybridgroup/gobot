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

	button := gobotGPIO.NewMakeyButton(firmata)
	button.Name = "button"
	button.Pin = "2"

	led := gobotGPIO.NewLed(firmata)
	led.Name = "led"
	led.Pin = "13"

	work := func() {
		gobot.On(button.Events["push"], func(data interface{}) {
			led.On()
		})

		gobot.On(button.Events["release"], func(data interface{}) {
			led.Off()
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmata},
		Devices:     []gobot.Device{button, led},
		Work:        work,
	}

	robot.Start()
}
