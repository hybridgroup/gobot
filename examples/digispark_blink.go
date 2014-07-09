package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/digispark"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func main() {
	gbot := gobot.NewGobot()
	digisparkAdaptor := digispark.NewDigisparkAdaptor("Digispark")
	led := gpio.NewLedDriver("led", digisparkAdaptor, "0")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{digisparkAdaptor},
		[]gobot.Device{led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
