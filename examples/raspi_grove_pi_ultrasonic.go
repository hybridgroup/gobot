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
	"gobot.io/x/gobot/v2/platforms/raspi"
)

const (
	ultrasonicPin = "4"
	delayMillisec = 10
)

func main() {
	r := raspi.NewAdaptor()
	gp := i2c.NewGrovePiDriver(r)

	work := func() {
		gobot.Every(1*time.Second, func() {
			if val, err := gp.UltrasonicRead(ultrasonicPin, delayMillisec); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Distance [cm]", val)
			}
		})
	}

	robot := gobot.NewRobot("ultrasonicBot",
		[]gobot.Connection{r},
		[]gobot.Device{gp},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
