//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/chip"
)

func main() {
	a := chip.NewAdaptor()
	blinkm := i2c.NewBlinkMDriver(a)

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
		[]gobot.Connection{a},
		[]gobot.Device{blinkm},
		work,
	)

	robot.Start()
}
