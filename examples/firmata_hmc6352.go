// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
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

	robot.Start()
}
