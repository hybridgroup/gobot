package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/drivers/i2c"
	"github.com/hybridgroup/gobot/platforms/firmata"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	hmc6352 := i2c.NewHMC6352Driver(firmataAdaptor)

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			heading, _ := hmc6352.Heading()
			fmt.Println("Heading", heading)
		})
	}

	robot := gobot.NewRobot("hmc6352Bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{hmc6352},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
