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
	led := gpio.NewLedDriver("led", digisparkAdaptor, "0")

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

	gbot.AddRobot(robot)

	gbot.Start()
}
