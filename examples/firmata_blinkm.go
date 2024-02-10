//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_blinkm.go /dev/ttyACM0
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
	blinkm := i2c.NewBlinkMDriver(firmataAdaptor)

	work := func() {
		gobot.Every(3*time.Second, func() {
			r := byte(gobot.Rand(255))
			g := byte(gobot.Rand(255))
			b := byte(gobot.Rand(255))
			if err := blinkm.Rgb(r, g, b); err != nil {
				fmt.Println(err)
			}
			color, err := blinkm.Color()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("color", color)
		})
	}

	robot := gobot.NewRobot("blinkmBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{blinkm},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
