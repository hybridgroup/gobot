//go:build libusb
// +build libusb

/*
Package digispark provides the Gobot adaptor for the Digispark ATTiny-based USB development board.

Installing:

This package requires installing `libusb`.
Then you can install the package with:

	Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Example:

	package main

	import (
		"time"

		"gobot.io/x/gobot/v2"
		"gobot.io/x/gobot/v2/drivers/gpio"
		"gobot.io/x/gobot/v2/platforms/digispark"
	)

	func main() {
		digisparkAdaptor := digispark.NewAdaptor()
		led := gpio.NewLedDriver(digisparkAdaptor, "0")

		work := func() {
			gobot.Every(1*time.Second, func() {
				if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
			})
		}

		robot := gobot.NewRobot("blinkBot",
			[]gobot.Connection{digisparkAdaptor},
			[]gobot.Device{led},
			work,
		)

		if err := robot.Start(); err != nil {
			panic(err)
		}
	}

For further information refer to digispark README:
https://github.com/hybridgroup/gobot/blob/release/platforms/digispark/README.md
*/
package digispark // import "gobot.io/x/gobot/v2/platforms/digispark"
