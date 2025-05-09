//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_rgb_led.go /dev/ttyACM0
*/

//nolint:gosec // ok here
package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	board := firmata.NewAdaptor(os.Args[1])
	led := gpio.NewRgbLedDriver(board, "3", "5", "6")

	work := func() {
		gobot.Every(1*time.Second, func() {
			r := uint8(gobot.Rand(255))
			g := uint8(gobot.Rand(255))
			b := uint8(gobot.Rand(255))
			if err := led.SetRGB(r, g, b); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("rgbBot",
		[]gobot.Connection{board},
		[]gobot.Device{led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
