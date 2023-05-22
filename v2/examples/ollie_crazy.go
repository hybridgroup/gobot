//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/ble"
	"gobot.io/x/gobot/v2/platforms/sphero/ollie"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	ollieBot := ollie.NewDriver(bleAdaptor)

	work := func() {
		ollieBot.SetRGB(255, 0, 255)
		gobot.Every(1*time.Second, func() {
			// Ollie performs 'crazy-ollie' trick
			ollieBot.SetRawMotorValues(ollie.Forward, uint8(255), ollie.Forward, uint8(255))
		})
	}

	robot := gobot.NewRobot("ollieBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ollieBot},
		work,
	)

	robot.Start()
}
