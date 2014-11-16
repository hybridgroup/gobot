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
	led := gpio.NewDirectPinDriver(beagleboneAdaptor, "led", "P8_10")
	button := gpio.NewDirectPinDriver(beagleboneAdaptor, "button", "P8_9")

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			val, _ := button.DigitalRead()
			if val == 1 {
				led.DigitalWrite(1)
			} else {
				led.DigitalWrite(0)
			}
		})
	}

	robot := gobot.NewRobot("pinBot",
		[]gobot.Connection{beagleboneAdaptor},
		[]gobot.Device{led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
