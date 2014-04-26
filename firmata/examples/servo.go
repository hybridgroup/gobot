package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-firmata"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {

	firmata := new(gobotFirmata.FirmataAdaptor)
	firmata.Name = "firmata"
	firmata.Port = "/dev/ttyACM0"

	servo := gobotGPIO.NewServo(firmata)
	servo.Name = "servo"
	servo.Pin = "3"

	work := func() {
		gobot.Every("1s", func() {
			i := uint8(gobot.Rand(180))
			fmt.Println("Turning", i)
			servo.Move(i)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmata},
		Devices:     []gobot.Device{servo},
		Work:        work,
	}

	robot.Start()
}
