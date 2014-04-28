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

	led := gpio.NewLedDriver(firmataAdaptor)
	led.Name = "led"
	led.Pin = "13"

	work := func() {
		gobot.Every("1s", func() {
			led.Toggle()
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmataAdaptor},
		Devices:     []gobot.Device{led},
		Work:        work,
	}

	robot.Start()
}
