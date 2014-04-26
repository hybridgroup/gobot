package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-beaglebone"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {
	beaglebone := new(gobotBeaglebone.Beaglebone)
	beaglebone.Name = "beaglebone"

	servo := gobotGPIO.NewServo(beaglebone)
	servo.Name = "servo"
	servo.Pin = "P9_14"

	work := func() {
		gobot.Every("1s", func() {
			i := uint8(gobot.Rand(180))
			fmt.Println("Turning", i)
			servo.Move(i)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{beaglebone},
		Devices:     []gobot.Device{servo},
		Work:        work,
	}

	robot.Start()
}
