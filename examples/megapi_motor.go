package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/megapi"
	"time"
)

func main() {
	gbot := gobot.NewGobot()

	// use "/dev/ttyUSB0" if connecting with USB cable
	// use "/dev/ttyAMA0" on devices older than Raspberry Pi 3 Model B
	megaPiAdaptor := megapi.NewMegaPiAdaptor("megapi", "/dev/ttyS0")
	motor := megapi.NewMotorDriver(megaPiAdaptor, "motor1", 1)

	work := func() {
		speed := int16(0)
		fadeAmount := int16(30)

		gobot.Every(100*time.Millisecond, func() {
			motor.Speed(speed)
			speed = speed + fadeAmount
			if speed == 0 || speed == 300 {
				fadeAmount = -fadeAmount
			}
		})
	}

	robot := gobot.NewRobot("megaPiBot",
		[]gobot.Connection{megaPiAdaptor},
		[]gobot.Device{motor},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
