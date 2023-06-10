/*
Package pebble contains the Gobot adaptor and driver for Pebble smart watch.

Installing:

It requires the 2.x iOS or Android app, and "watchbot" app (https://gobot.io/x/watchbot)
installed on Pebble watch. Then install running:

	Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Example:

Before running the example, make sure configuration settings match with your program. In the example, api host is your computer IP, robot name is 'pebble' and robot api port is 8080

	package main

	import (
		"fmt"

		"gobot.io/x/gobot/v2"
		"gobot.io/x/gobot/v2/api"
		"gobot.io/x/gobot/v2/platforms/pebble"
	)

	func main() {
		master := gobot.NewMaster()
		api.NewAPI(master).Start()

		pebbleAdaptor := pebble.NewAdaptor()
		watch := pebble.NewDriver(pebbleAdaptor)

		work := func() {
			watch.SendNotification("Hello Pebble!")
			watch.On(watch.Event("button"), func(data interface{}) {
				fmt.Println("Button pushed: " + data.(string))
			})

			watch.On(watch.Event("tap"), func(data interface{}) {
				fmt.Println("Tap event detected")
			})
		}

		robot := gobot.NewRobot("pebble",
			[]gobot.Connection{pebbleAdaptor},
			[]gobot.Device{watch},
			work,
		)

		master.AddRobot(robot)

		master.Start()
	}

For more information refer to the pebble README:
https://github.com/hybridgroup/gobot/blob/master/platforms/pebble/README.md
*/
package pebble // import "gobot.io/x/gobot/v2/platforms/pebble"
