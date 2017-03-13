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
	lidar := i2c.NewLIDARLiteDriver(firmataAdaptor)

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			distance, _ := lidar.Distance()
			fmt.Println("Distance", distance)
		})
	}

	robot := gobot.NewRobot("lidarbot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{lidar},
		work,
	)

	robot.Start()
}
