/*
Package particle provides the Gobot adaptor for the Particle Photon and Electron.

Installing:

	Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Example:

	package main

	import (
		"time"

		"gobot.io/x/gobot/v2"
		"gobot.io/x/gobot/v2/drivers/gpio"
		"gobot.io/x/gobot/v2/platforms/particle"
	)

	func main() {
		core := paticle.NewAdaptor("device_id", "access_token")
		led := gpio.NewLedDriver(core, "D7")

		work := func() {
			gobot.Every(1*time.Second, func() {
				if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
			})
		}

		robot := gobot.NewRobot("particle",
			[]gobot.Connection{core},
			[]gobot.Device{led},
			work,
		)

		if err := robot.Start(); err != nil {
			panic(err)
		}
	}

For further information refer to Particle readme:
https://github.com/hybridgroup/gobot/blob/release/platforms/particle/README.md
*/
package particle // import "gobot.io/x/gobot/v2/platforms/particle"
