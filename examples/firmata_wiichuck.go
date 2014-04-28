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

	wiichuck := i2c.NewWiichuckDriver(firmataAdaptor)
	wiichuck.Name = "wiichuck"

	work := func() {
		gobot.On(wiichuck.Events["joystick"], func(data interface{}) {
			fmt.Println("joystick")
		})

		gobot.On(wiichuck.Events["c_button"], func(data interface{}) {
			fmt.Println("c")
		})

		gobot.On(wiichuck.Events["z_button"], func(data interface{}) {
			fmt.Println("z")
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmataAdaptor},
		Devices:     []gobot.Device{wiichuck},
		Work:        work,
	}

	robot.Start()
}
