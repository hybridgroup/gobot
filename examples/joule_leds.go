// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/intel-iot/joule"
)

func main() {
	e := joule.NewAdaptor()
	led0 := gpio.NewLedDriver(e, "GP100")
	led1 := gpio.NewLedDriver(e, "GP101")
	led2 := gpio.NewLedDriver(e, "GP102")
	led3 := gpio.NewLedDriver(e, "GP103")

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
		gobot.Every(4*time.Second, func() {
			led2.Toggle()
		})
		gobot.Every(8*time.Second, func() {
			led3.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{e},
		[]gobot.Device{led0, led1, led2, led3},
		work,
	)

	robot.Start()
}
