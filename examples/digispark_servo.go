package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/digispark"
	"github.com/hybridgroup/gobot/gpio"
)

func main() {
	digisparkAdaptor := digispark.NewDigisparkAdaptor()
	digisparkAdaptor.Name = "digispark"

	servo := gpio.NewServoDriver(digisparkAdaptor)
	servo.Name = "servo"
	servo.Pin = "0"

	work := func() {
		gobot.Every("1s", func() {
			i := uint8(gobot.Rand(180))
			fmt.Println("Turning", i)
			servo.Move(i)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{digisparkAdaptor},
		Devices:     []gobot.Device{servo},
		Work:        work,
	}

	robot.Start()
}
