/*
Package leap provides the Gobot adaptor and driver for the Leap Motion.

Installing:

* First install the [Leap Motion Software](https://www.leapmotion.com/setup).
* Then install the package:

	go get github.com/hybridgroup/gobot/platforms/leap

Example:

	package main

	import (
		"fmt"

		"github.com/hybridgroup/gobot"
		"github.com/hybridgroup/gobot/platforms/leap"
	)

	func main() {
		gbot := gobot.NewGobot()

		leapMotionAdaptor := leap.NewLeapMotionAdaptor("leap", "127.0.0.1:6437")
		l := leap.NewLeapMotionDriver(leapMotionAdaptor, "leap")

		work := func() {
			gobot.On(l.Event("message"), func(data interface{}) {
				fmt.Println(data.(leap.Frame))
			})
		}

		robot := gobot.NewRobot("leapBot",
			[]gobot.Connection{leapMotionAdaptor},
			[]gobot.Device{l},
			work,
		)

		gbot.AddRobot(robot)

		gbot.Start()
	}

For more information refer to the leap README:
https://github.com/hybridgroup/gobot/blob/master/platforms/leap/README.md
*/
package leap
