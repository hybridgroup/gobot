package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/digispark"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	digisparkAdaptor := digispark.NewDigisparkAdaptor("digispark")
	servo := gpio.NewServoDriver(digisparkAdaptor, "servo", "0")

	work := func() {
		gobot.Every(1*time.Second, func() {
			i := uint8(gobot.Rand(180))
			fmt.Println("Turning", i)
			servo.Move(i)
		})
	}

	robot := gobot.NewRobot("servoBot",
		[]gobot.Connection{digisparkAdaptor},
		[]gobot.Device{servo},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
