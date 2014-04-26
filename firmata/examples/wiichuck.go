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

	wiichuck := gobotI2C.NewWiichuck(firmata)
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
		Connections: []gobot.Connection{firmata},
		Devices:     []gobot.Device{wiichuck},
		Work:        work,
	}

	robot.Start()
}
