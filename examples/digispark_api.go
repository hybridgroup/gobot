package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/digispark"
	"github.com/hybridgroup/gobot/gpio"
)

func main() {
	master := gobot.NewMaster()
	gobot.StartApi(master)

	digisparkAdaptor := digispark.NewDigisparkAdaptor()
	digisparkAdaptor.Name = "Digispark"

	led := gpio.NewLed(digisparkAdaptor)
	led.Name = "led"
	led.Pin = "0"

	master.Robots = append(master.Robots, &gobot.Robot{
		Name:        "digispark",
		Connections: []gobot.Connection{digisparkAdaptor},
		Devices:     []gobot.Device{led},
	})

	master.Start()
}
