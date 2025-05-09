/*
Package ardrone provides the Gobot adaptor and driver for the Parrot Ardrone.

Installing:

	Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Example:

	package main

	import (
		"time"

		"gobot.io/x/gobot/v2"
		"gobot.io/x/gobot/v2/platforms/parrot/ardrone"
	)

	func main() {
		ardroneAdaptor := ardrone.NewAdaptor()
		drone := ardrone.NewDriver(ardroneAdaptor)

		work := func() {
			drone.TakeOff()
			drone.On(drone.Event("flying"), func(data interface{}) {
				gobot.After(3*time.Second, func() {
					drone.Land()
				})
			})
		}

		robot := gobot.NewRobot("drone",
			[]gobot.Connection{ardroneAdaptor},
			[]gobot.Device{drone},
			work,
		)

		if err := robot.Start(); err != nil {
			panic(err)
		}
	}

For more information refer to the ardrone README:
https://github.com/hybridgroup/gobot/blob/release/platforms/parrot/ardrone/README.md
*/
package ardrone // import "gobot.io/x/gobot/v2/platforms/parrot/ardrone"
