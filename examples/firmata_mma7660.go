//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_mma7660.go /dev/ttyACM0
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

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
