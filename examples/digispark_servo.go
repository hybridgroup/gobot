package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/digispark"
	"github.com/hybridgroup/gobot/gpio"
	"time"
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

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("servoBot", []gobot.Connection{digisparkAdaptor}, []gobot.Device{servo}, work))
	gbot.Start()
}
