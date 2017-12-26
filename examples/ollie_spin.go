// +build example
//
// Do not build by default.

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/sphero/ollie"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	ollieBot := ollie.NewDriver(bleAdaptor)

	work := func() {
		ollieBot.SetRGB(255, 0, 255)
		gobot.Every(1*time.Second, func() {
			// Ollie performs 360 spin trick
			ollieBot.SetRawMotorValues(ollie.Forward, uint8(255), ollie.Reverse, uint8(255))
		})
	}

	robot := gobot.NewRobot("ollieBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ollieBot},
		work,
	)

	robot.Start()
}
