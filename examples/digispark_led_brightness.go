// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/digispark"
)

func main() {
	digisparkAdaptor := digispark.NewAdaptor()
	led := gpio.NewLedDriver(digisparkAdaptor, "0")

	work := func() {
		brightness := uint8(0)
		fadeAmount := uint8(15)

		gobot.Every(100*time.Millisecond, func() {
			led.Brightness(brightness)
			brightness = brightness + fadeAmount
			if brightness == 0 || brightness == 255 {
				fadeAmount = -fadeAmount
			}
		})
	}

	robot := gobot.NewRobot("pwmBot",
		[]gobot.Connection{digisparkAdaptor},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
