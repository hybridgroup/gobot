/*
Package sphero provides the Gobot adaptor and driver for the Sphero.

Installing:

	go get github.com/hybridgroup/gobot/platforms/sphero

Example:

	package main

	import (
		"fmt"
		"time"

		"github.com/hybridgroup/gobot"
		"github.com/hybridgroup/gobot/platforms/sphero"
	)

	func main() {
		gbot := gobot.NewGobot()

		adaptor := sphero.NewSpheroAdaptor("sphero", "/dev/rfcomm0")
		driver := sphero.NewSpheroDriver(adaptor, "sphero")

		work := func() {
			gobot.Every(3*time.Second, func() {
				driver.Roll(30, uint16(gobot.Rand(360)))
			})
		}

		robot := gobot.NewRobot("sphero",
			[]gobot.Connection{adaptor},
			[]gobot.Device{driver},
			work,
		)

		gbot.AddRobot(robot)

		gbot.Start()
	}

For futher information refer to sphero readme:
https://github.com/hybridgroup/gobot/blob/master/platforms/sphero/README.md
*/
package sphero
