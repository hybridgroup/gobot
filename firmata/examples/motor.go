package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-firmata"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {
	firmata := new(gobotFirmata.FirmataAdaptor)
	firmata.Name = "firmata"
	firmata.Port = "/dev/ttyACM0"

	motor := gobotGPIO.NewMotor(firmata)
	motor.Name = "motor"
	motor.SpeedPin = "3"

	work := func() {
		speed := byte(0)
		fade_amount := byte(15)

		gobot.Every("0.1s", func() {
			motor.Speed(speed)
			speed = speed + fade_amount
			if speed == 0 || speed == 255 {
				fade_amount = -fade_amount
			}
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmata},
		Devices:     []gobot.Device{motor},
		Work:        work,
	}

	robot.Start()
}
