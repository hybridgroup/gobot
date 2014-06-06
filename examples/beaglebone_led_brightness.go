package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func main() {
	gbot := gobot.NewGobot()
	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
	led := gpio.NewLedDriver(beagleboneAdaptor, "led", "P9_14")

	work := func() {
		brightness := uint8(0)
		fade_amount := uint8(5)

		gobot.Every(100*time.Millisecond, func() {
			led.Brightness(brightness)
			brightness = brightness + fade_amount
			if brightness == 0 || brightness == 255 {
				fade_amount = -fade_amount
			}
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("pwmBot", []gobot.Connection{beagleboneAdaptor}, []gobot.Device{led}, work))
	gbot.Start()
}
