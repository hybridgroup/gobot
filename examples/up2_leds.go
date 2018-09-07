// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/upboard/up2"
)

func main() {
	b := up2.NewAdaptor()
	red := gpio.NewLedDriver(b, "red")
	blue := gpio.NewLedDriver(b, "blue")
	green := gpio.NewLedDriver(b, "green")
	yellow := gpio.NewLedDriver(b, "yellow")

	work := func() {
		red.Off()
		blue.Off()
		green.Off()
		yellow.Off()

		gobot.Every(1*time.Second, func() {
			red.Toggle()
		})
		gobot.Every(2*time.Second, func() {
			green.Toggle()
		})
		gobot.Every(4*time.Second, func() {
			yellow.Toggle()
		})
		gobot.Every(8*time.Second, func() {
			blue.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{b},
		[]gobot.Device{red, blue, green, yellow},
		work,
	)

	robot.Start()
}
