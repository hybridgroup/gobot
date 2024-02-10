//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/sphero"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])
	ollieBot := sphero.NewOllieDriver(bleAdaptor)

	work := func() {
		ollieBot.SetRGB(255, 0, 255)
		gobot.Every(1*time.Second, func() {
			// Ollie performs 'crazy-ollie' trick
			ollieBot.SetRawMotorValues(sphero.Forward, uint8(255), sphero.Forward, uint8(255))
		})
	}

	robot := gobot.NewRobot("ollieBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ollieBot},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
