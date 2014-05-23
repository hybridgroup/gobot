package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/digispark"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	master := gobot.NewGobot()
	api.NewApi(master).Start()

	digisparkAdaptor := digispark.NewDigisparkAdaptor("Digispark")
	led := gpio.NewLed(digisparkAdaptor, "led", "0")

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("digispark", []gobot.Connection{digisparkAdaptor}, []gobot.Device{led}, nil))
	gbot.Start()
}
