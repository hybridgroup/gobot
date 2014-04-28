package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/firmata"
	"github.com/hybridgroup/gobot/gpio"
)

func main() {
	firmataAdaptor := firmata.NewFirmataAdaptor()
	firmataAdaptor.Name = "firmata"
	firmataAdaptor.Port = "/dev/ttyACM0"

	button := gpio.NewButtonDriver(firmataAdaptor)
	button.Name = "button"
	button.Pin = "2"

	led := gpio.NewLedDriver(firmataAdaptor)
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
		Connections: []gobot.Connection{firmataAdaptor},
		Devices:     []gobot.Device{button, led},
		Work:        work,
	}

	robot.Start()
}
