package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/particle"
)

func main() {
	core := particle.NewAdaptor("device_id", "access_token")
	led := gpio.NewLedDriver(core, "A1")

	work := func() {
		brightness := uint8(0)
		fadeAmount := uint8(25)

		gobot.Every(500*time.Millisecond, func() {
			led.Brightness(brightness)
			brightness = brightness + fadeAmount
			if brightness == 0 || brightness == 255 {
				fadeAmount = -fadeAmount
			}
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
