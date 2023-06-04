/*
Package beaglebone provides the Gobot adaptor for the Beaglebone Black/Green, as well as a
separate Adaptor for the PocketBeagle.

Installing:

	Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Example:

	package main

	import (
		"time"

		"gobot.io/x/gobot/v2"
		"gobot.io/x/gobot/v2/drivers/gpio"
		"gobot.io/x/gobot/v2/platforms/beaglebone"
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
package beaglebone // import "gobot.io/x/gobot/v2/platforms/beaglebone"
