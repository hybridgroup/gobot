package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/drivers/gpio"
	"github.com/hybridgroup/gobot/platforms/particle"
)

func main() {
	gbot := gobot.NewMaster()

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

	gbot.AddRobot(robot)

	gbot.Start()
}
