//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_servo.go /dev/ttyACM0
*/

//nolint:gosec // ok here
package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	servo := gpio.NewServoDriver(firmataAdaptor, "3")

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
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{servo},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
