// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	pin := gpio.NewDirectPinDriver(firmataAdaptor, "13")

	work := func() {
		level := byte(1)

		gobot.Every(1*time.Second, func() {
			pin.DigitalWrite(level)
			if level == 1 {
				level = 0
			} else {
				level = 1
			}
		})
	}

	robot := gobot.NewRobot("pinBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{pin},
		work,
	)

	robot.Start()
}
