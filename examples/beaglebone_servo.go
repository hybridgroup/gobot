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
	"gobot.io/x/gobot/v2/platforms/beaglebone"
)

func main() {
	beagleboneAdaptor := beaglebone.NewAdaptor()
	servo := gpio.NewServoDriver(beagleboneAdaptor, "P9_14")

	work := func() {
		gobot.Every(1*time.Second, func() {
			i := uint8(gobot.Rand(180))
			fmt.Println("Turning", i)
			if err := servo.Move(i); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("servoBot",
		[]gobot.Connection{beagleboneAdaptor},
		[]gobot.Device{servo},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
