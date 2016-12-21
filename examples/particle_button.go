package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/particle"
)

func main() {
	core := particle.NewAdaptor("device_id", "access_token")
	led := gpio.NewLedDriver(core, "D7")
	button := gpio.NewButtonDriver(core, "D5")

	work := func() {
		button.On(button.Event("push"), func(data interface{}) {
			led.On()
		})

		button.On(button.Event("release"), func(data interface{}) {
			led.Off()
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		[]gobot.Device{button, led},
		work,
	)

	robot.Start()
}
