package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func main() {
	gbot := gobot.NewGobot()
	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/ttyACM0")
	motor := gpio.NewMotorDriver(firmataAdaptor, "motor", "3")

	work := func() {
		speed := byte(0)
		fade_amount := byte(15)

		gobot.Every(100*time.Millisecond, func() {
			motor.Speed(speed)
			speed = speed + fade_amount
			if speed == 0 || speed == 255 {
				fade_amount = -fade_amount
			}
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("motorBot", []gobot.Connection{firmataAdaptor}, []gobot.Device{motor}, work))

	gbot.Start()
}
