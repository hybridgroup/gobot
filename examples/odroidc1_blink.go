package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/dhart-alldigital/gobot/platforms/odroid/c1"
)

func main() {
	gbot := gobot.NewGobot()

	r := c1.NewODroidC1Adaptor("c1")
	led := gpio.NewLedDriver(r, "led", "11")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
