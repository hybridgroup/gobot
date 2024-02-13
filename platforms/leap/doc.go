/*
Package leap provides the Gobot adaptor and driver for the Leap Motion.

Installing:

	Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

	Install the [Leap Motion Software](https://www.leapmotion.com/setup).

Example:

	package main

	import (
		"fmt"

		"gobot.io/x/gobot/v2"
		"gobot.io/x/gobot/v2/platforms/leap"
	)

	func main() {
		leapMotionAdaptor := leap.NewAdaptor("127.0.0.1:6437")
		l := leap.NewDriver(leapMotionAdaptor)

		work := func() {
			l.On(l.Event("message"), func(data interface{}) {
				fmt.Println(data.(leap.Frame))
			})
		}

		robot := gobot.NewRobot("leapBot",
			[]gobot.Connection{leapMotionAdaptor},
			[]gobot.Device{l},
			work,
		)

		if err := robot.Start(); err != nil {
			panic(err)
		}
	}

For more information refer to the leap README:
https://github.com/hybridgroup/gobot/blob/release/platforms/leap/README.md
*/
package leap // import "gobot.io/x/gobot/v2/platforms/leap"
