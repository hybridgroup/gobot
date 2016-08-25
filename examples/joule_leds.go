package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/joule"
)

func main() {
	gbot := gobot.NewGobot()

	e := joule.NewJouleAdaptor("joule")
	led0 := gpio.NewLedDriver(e, "led", "100")
	led1 := gpio.NewLedDriver(e, "led", "101")
	led2 := gpio.NewLedDriver(e, "led", "102")
	led3 := gpio.NewLedDriver(e, "led", "103")

	work := func() {
		led0.Off()
		led1.Off()
		led2.Off()
		led3.Off()

		gobot.Every(1*time.Second, func() {
			led0.Toggle()
		})
		gobot.Every(2*time.Second, func() {
			led1.Toggle()
		})
		gobot.Every(2*time.Second, func() {
			led2.Toggle()
		})
		gobot.Every(3*time.Second, func() {
			led3.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{e},
		[]gobot.Device{led0, led1, led2, led3},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
