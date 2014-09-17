package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
	led := gpio.NewLedDriver(beagleboneAdaptor, "led", "usr0")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{beagleboneAdaptor},
		[]gobot.Device{led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
