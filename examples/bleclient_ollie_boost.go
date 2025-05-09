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
		head := uint16(90)
		ollieBot.SetRGB(255, 0, 0)
		ollieBot.Boost(true)
		gobot.Every(1*time.Second, func() {
			ollieBot.Roll(0, head)
			time.Sleep(1 * time.Second)
			head += 90
			head = head % 360
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
