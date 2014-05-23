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
	led := gpio.NewLedDriver(digisparkAdaptor, "led", "0")

	work := func() {
		gobot.Every(0.5*time.Second, func() {
			led.Toggle()
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("blinkBot", []gobot.Connection{digisparkAdaptor}, []gobot.Device{led}, work))
	gbot.Start()
}
