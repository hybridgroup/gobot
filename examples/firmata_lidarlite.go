//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_lidarlite.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
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

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
