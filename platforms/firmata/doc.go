/*
Package firmata provides the Gobot adaptor for microcontrollers that support the Firmata protocol.

Installing:

	go get -d -u github.com/hybridgroup/gobot/... && go get github.com/hybridgroup/gobot/platforms/firmata

Example:

	package main

	import (
		"time"

		"github.com/hybridgroup/gobot"
		"github.com/hybridgroup/gobot/drivers/gpio"
		"github.com/hybridgroup/gobot/platforms/firmata"
	)

	func main() {
		firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
		led := gpio.NewLedDriver(firmataAdaptor, "13")

		work := func() {
			gobot.Every(1*time.Second, func() {
				led.Toggle()
			})
		}

		robot := gobot.NewRobot("bot",
			[]gobot.Connection{firmataAdaptor},
			[]gobot.Device{led},
			work,
		)

		robot.Start()
	}

For further information refer to firmata readme:
https://github.com/hybridgroup/gobot/blob/master/platforms/firmata/README.md
*/
package firmata
