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
	mma7660 := i2c.NewMMA7660Driver(firmataAdaptor)

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

	robot.Start()
}
