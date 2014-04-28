package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/firmata"
	"github.com/hybridgroup/gobot/i2c"
)

func main() {
	firmataAdaptor := firmata.NewFirmataAdaptor()
	firmataAdaptor.Name = "firmata"
	firmataAdaptor.Port = "/dev/ttyACM0"

	blinkm := i2c.NewBlinkMDriver(firmataAdaptor)
	blinkm.Name = "blinkm"

	work := func() {
		gobot.Every("3s", func() {
			blinkm.Rgb(byte(gobot.Rand(255)), byte(gobot.Rand(255)), byte(gobot.Rand(255)))
			fmt.Println("color", blinkm.Color())
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmataAdaptor},
		Devices:     []gobot.Device{blinkm},
		Work:        work,
	}

	robot.Start()
}
