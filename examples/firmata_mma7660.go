package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/i2c"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/ttyACM0")
	mma7660 := i2c.NewMMA7660Driver(firmataAdaptor, "mma7660")

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			if x, y, z, err := mma7660.XYZ(); err == nil {
				fmt.Println(x, y, z)
				fmt.Println(mma7660.Acceleration(x, y, z))
			} else {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("mma76602Bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{mma7660},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
