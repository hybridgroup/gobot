package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func main() {
	gbot := gobot.NewGobot()
	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "led", "3")

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
		gobot.NewRobot("pwmBot", []gobot.Connection{firmataAdaptor}, []gobot.Device{led}, work))
	gbot.Start()

}
