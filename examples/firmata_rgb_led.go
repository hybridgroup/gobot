// +build example
//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_rgb_led.go /dev/ttyACM0
*/

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	board := firmata.NewAdaptor(os.Args[1])
	led := gpio.NewRgbLedDriver(board, "3", "5", "6")

	work := func() {
		gobot.Every(1*time.Second, func() {
			r := uint8(gobot.Rand(255))
			g := uint8(gobot.Rand(255))
			b := uint8(gobot.Rand(255))
			led.SetRGB(r, g, b)
		})
	}

	robot := gobot.NewRobot("rgbBot",
		[]gobot.Connection{board},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
