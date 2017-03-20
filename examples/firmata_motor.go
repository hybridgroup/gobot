// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	motor := gpio.NewMotorDriver(firmataAdaptor, "3")

	work := func() {
		speed := byte(0)
		fadeAmount := byte(15)

		gobot.Every(100*time.Millisecond, func() {
			motor.Speed(speed)
			speed = speed + fadeAmount
			if speed == 0 || speed == 255 {
				fadeAmount = -fadeAmount
			}
		})
	}

	robot := gobot.NewRobot("motorBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{motor},
		work,
	)

	robot.Start()
}
