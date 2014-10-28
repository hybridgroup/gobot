/*
Package spark provides the Gobot adaptor for the Spark Core.

Installing:

	go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/gobot/platforms/spark

Example:

	package main

	import (
		"time"

		"github.com/hybridgroup/gobot"
		"github.com/hybridgroup/gobot/platforms/gpio"
		"github.com/hybridgroup/gobot/platforms/spark"
	)

	func main() {
		gbot := gobot.NewGobot()

		sparkCore := spark.NewSparkCoreAdaptor("spark", "device_id", "access_token")
		led := gpio.NewLedDriver(sparkCore, "led", "D7")

		work := func() {
			gobot.Every(1*time.Second, func() {
				led.Toggle()
			})
		}

		robot := gobot.NewRobot("spark",
			[]gobot.Connection{sparkCore},
			[]gobot.Device{led},
			work,
		)

		gbot.AddRobot(robot)

		gbot.Start()
	}

For further information refer to spark readme:
https://github.com/hybridgroup/gobot/blob/master/platforms/spark/README.md
*/
package spark
