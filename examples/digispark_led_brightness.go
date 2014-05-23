package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/digispark"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func main() {
	gbot := gobot.NewGobot()

	digisparkAdaptor := digispark.NewDigisparkAdaptor("digispark")
	led := gpio.NewLedDriver(digisparkAdaptor, "led", "0")

	work := func() {
		brightness := uint8(0)
		fade_amount := uint8(15)

		gobot.Every(0.1*time.Second, func() {
			led.Brightness(brightness)
			brightness = brightness + fade_amount
			if brightness == 0 || brightness == 255 {
				fade_amount = -fade_amount
			}
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("pwmBot", []gobot.Connection{digisparkAdaptor}, []gobot.Device{led}, work))
	gbot.Start()
}
