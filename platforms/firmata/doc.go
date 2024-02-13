/*
Package firmata provides the Gobot adaptor for microcontrollers that support the Firmata protocol.

Installing:

	Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Example:

	package main

	import (
		"time"

		"gobot.io/x/gobot/v2"
		"gobot.io/x/gobot/v2/drivers/gpio"
		"gobot.io/x/gobot/v2/platforms/firmata"
	)

	func main() {
		firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
		led := gpio.NewLedDriver(firmataAdaptor, "13")

		work := func() {
			gobot.Every(1*time.Second, func() {
				if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
			})
		}

		robot := gobot.NewRobot("bot",
			[]gobot.Connection{firmataAdaptor},
			[]gobot.Device{led},
			work,
		)

		if err := robot.Start(); err != nil {
			panic(err)
		}
	}

For further information refer to firmata readme:
https://github.com/hybridgroup/gobot/blob/release/platforms/firmata/README.md
*/
package firmata // import "gobot.io/x/gobot/v2/platforms/firmata"
