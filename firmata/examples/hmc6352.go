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

	hmc6352 := gobotI2C.NewHMC6352(firmata)
	hmc6352.Name = "hmc6352"

	work := func() {
		gobot.Every("0.1s", func() {
			fmt.Println("Heading", hmc6352.Heading)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmata},
		Devices:     []gobot.Device{hmc6352},
		Work:        work,
	}

	robot.Start()
}
