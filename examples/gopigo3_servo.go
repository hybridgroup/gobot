//go:build example
// +build example

//
// Do not build by default.

//nolint:gosec // ok here
package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/dexter/gopigo3"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	raspiAdaptor := raspi.NewAdaptor()
	gpg3 := gopigo3.NewDriver(raspiAdaptor)
	servo := gpio.NewServoDriver(gpg3, "SERVO_1")

	work := func() {
		gobot.Every(1*time.Second, func() {
			i := uint8(gobot.Rand(180))
			fmt.Println("Turning", i)
			if err := servo.Move(i); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("gopigo3servo",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{gpg3, servo},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
