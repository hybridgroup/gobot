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

	gbot.AddRobot(robot)

	gbot.Start()
}
