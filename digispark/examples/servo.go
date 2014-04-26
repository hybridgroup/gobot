package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-digispark"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {

	digispark := new(gobotDigispark.DigisparkAdaptor)
	digispark.Name = "digispark"

	servo := gobotGPIO.NewServo(digispark)
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
		Connections: []gobot.Connection{digispark},
		Devices:     []gobot.Device{servo},
		Work:        work,
	}

	robot.Start()
}
