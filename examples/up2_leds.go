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
	red := gpio.NewLedDriver(b, up2.LEDRed)
	blue := gpio.NewLedDriver(b, up2.LEDBlue)
	green := gpio.NewLedDriver(b, up2.LEDGreen)
	yellow := gpio.NewLedDriver(b, up2.LEDYellow)

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
