/*
Package beaglebone provides the Gobot adaptor for the Beaglebone Black.

Installing:

	go get gobot.io/x/gobot/platforms/beaglebone

Example:

	package main

	import (
		"time"

		"gobot.io/x/gobot"
		"gobot.io/x/gobot/drivers/gpio"
		"gobot.io/x/gobot/platforms/beaglebone"
	)

	func main() {
		beagleboneAdaptor := beaglebone.NewAdaptor()
		led := gpio.NewLedDriver(beagleboneAdaptor, "P9_12")

		work := func() {
			gobot.Every(1*time.Second, func() {
				led.Toggle()
			})
		}

		robot := gobot.NewRobot("blinkBot",
			[]gobot.Connection{beagleboneAdaptor},
			[]gobot.Device{led},
			work,
		)

		robot.Start()
	}

For more information refer to the beaglebone README:
https://github.com/hybridgroup/gobot/blob/master/platforms/beaglebone/README.md
*/
package beaglebone // import "gobot.io/x/gobot/platforms/beaglebone"
