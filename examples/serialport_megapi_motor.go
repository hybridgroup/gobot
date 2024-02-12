//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/serial/megapi"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
	// use "/dev/ttyUSB0" if connecting with USB cable
	// use "/dev/ttyAMA0" on devices older than Raspberry Pi 3 Model B
	adaptor := serialport.NewAdaptor("/dev/ttyS0", serialport.WithName("MegaPi"))
	motor := megapi.NewMotorDriver(adaptor, 1)

	work := func() {
		speed := int16(0)
		fadeAmount := int16(30)

		gobot.Every(100*time.Millisecond, func() {
			if err := motor.Speed(speed); err != nil {
				fmt.Println(err)
			}
			speed = speed + fadeAmount
			if speed == 0 || speed == 300 {
				fadeAmount = -fadeAmount
			}
		})
	}

	robot := gobot.NewRobot("megaPiBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{motor},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
