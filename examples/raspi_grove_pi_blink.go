//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	gp := i2c.NewGrovePiDriver(r)
	led := gpio.NewLedDriver(gp, "D2")

	work := func() {
		gobot.Every(1*time.Second, func() {
			if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{gp, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
