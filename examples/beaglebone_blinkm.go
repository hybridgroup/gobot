// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/beaglebone"
)

func main() {
	beagleboneAdaptor := beaglebone.NewAdaptor()
	blinkm := i2c.NewBlinkMDriver(beagleboneAdaptor)

	work := func() {
		gobot.Every(1*time.Second, func() {
			r := byte(gobot.Rand(255))
			g := byte(gobot.Rand(255))
			b := byte(gobot.Rand(255))
			blinkm.Rgb(r, g, b)
			color, _ := blinkm.Color()
			fmt.Println("color", color)
		})
	}

	robot := gobot.NewRobot("blinkmBot",
		[]gobot.Connection{beagleboneAdaptor},
		[]gobot.Device{blinkm},
		work,
	)

	robot.Start()
}
