/*
Package spark provides the Gobot adaptor for the Spark Core.

Installing:

	go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/gobot/platforms/spark

Example:

	package main

	import (
		"time"

		"github.com/hybridgroup/gobot"
		"github.com/hybridgroup/gobot/drivers/gpio"
		"github.com/hybridgroup/gobot/platforms/particle"
	)

	func main() {
		core := paticle.NewAdaptor("device_id", "access_token")
		led := gpio.NewLedDriver(core, "D7")

		work := func() {
			gobot.Every(1*time.Second, func() {
				led.Toggle()
			})
		}

		robot := gobot.NewRobot("spark",
			[]gobot.Connection{core},
			[]gobot.Device{led},
			work,
		)

		robot.Start()
	}

For further information refer to Particle readme:
https://github.com/hybridgroup/gobot/blob/master/platforms/particle/README.md
*/
package particle
