//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_hmc6352.go /dev/ttyACM0
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

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
