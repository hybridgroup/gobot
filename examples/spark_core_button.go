package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/spark"
)

func main() {
	gbot := gobot.NewGobot()

	sparkCore := spark.NewSparkCoreAdaptor("spark", "device_id", "access_token")
	led := gpio.NewLedDriver(sparkCore, "led", "D7")
	button := gpio.NewButtonDriver(sparkCore, "button", "D5")

	work := func() {
		gobot.On(button.Event("push"), func(data interface{}) {
			led.On()
		})

		gobot.On(button.Event("release"), func(data interface{}) {
			led.Off()
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{sparkCore},
		[]gobot.Device{button, led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
