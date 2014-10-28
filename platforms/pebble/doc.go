/*
Package pebble contains the Gobot adaptor and driver for Pebble smart watch.

Installing:

It requires the 2.x iOS or Android app, and "watchbot" app (https://github.com/hybridgroup/watchbot)
installed on Pebble watch. Then install running:

	go get github.com/hybridgroup/gobot/platforms/pebble

Example:

Before running the example, make sure configuration settings match with your program. In the example, api host is your computer IP, robot name is 'pebble' and robot api port is 8080

	package main

	import (
		"fmt"

		"github.com/hybridgroup/gobot"
		"github.com/hybridgroup/gobot/api"
		"github.com/hybridgroup/gobot/platforms/pebble"
	)

	func main() {
		gbot := gobot.NewGobot()
		api.NewAPI(gbot).Start()

		pebbleAdaptor := pebble.NewPebbleAdaptor("pebble")
		pebbleDriver := pebble.NewPebbleDriver(pebbleAdaptor, "pebble")

		work := func() {
			pebbleDriver.SendNotification("Hello Pebble!")
			gobot.On(pebbleDriver.Event("button"), func(data interface{}) {
				fmt.Println("Button pushed: " + data.(string))
			})

			gobot.On(pebbleDriver.Event("tap"), func(data interface{}) {
				fmt.Println("Tap event detected")
			})
		}

		robot := gobot.NewRobot("pebble",
			[]gobot.Connection{pebbleAdaptor},
			[]gobot.Device{pebbleDriver},
			work,
		)

		gbot.AddRobot(robot)

		gbot.Start()
	}

For more information refer to the pebble README:
https://github.com/hybridgroup/gobot/blob/master/platforms/pebble/README.md
*/
package pebble
