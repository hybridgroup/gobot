package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/i2c"
)

func main() {
	gbot := gobot.NewGobot()
	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/ttyACM0")
	wiichuck := i2c.NewWiichuckDriver(firmataAdaptor, "wiichuck")

	work := func() {
		gobot.On(wiichuck.Events["joystick"], func(data interface{}) {
			fmt.Println("joystick", data)
		})

		gobot.On(wiichuck.Events["c_button"], func(data interface{}) {
			fmt.Println("c")
		})

		gobot.On(wiichuck.Events["z_button"], func(data interface{}) {
			fmt.Println("z")
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("chuck", []gobot.Connection{firmataAdaptor}, []gobot.Device{wiichuck}, work))

	gbot.Start()
}
