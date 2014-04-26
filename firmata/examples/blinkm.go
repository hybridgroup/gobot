package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-firmata"
	"github.com/hybridgroup/gobot-i2c"
)

func main() {
	firmata := new(gobotFirmata.FirmataAdaptor)
	firmata.Name = "firmata"
	firmata.Port = "/dev/ttyACM0"

	blinkm := gobotI2C.NewBlinkM(firmata)
	blinkm.Name = "blinkm"

	work := func() {
		gobot.Every("3s", func() {
			blinkm.Rgb(byte(gobot.Rand(255)), byte(gobot.Rand(255)), byte(gobot.Rand(255)))
			fmt.Println("color", blinkm.Color())
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmata},
		Devices:     []gobot.Device{blinkm},
		Work:        work,
	}

	robot.Start()
}
