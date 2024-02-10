//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_direct_pin.go /dev/ttyACM0
*/

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
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	pin := gpio.NewDirectPinDriver(firmataAdaptor, "13")

	work := func() {
		level := byte(1)

		gobot.Every(1*time.Second, func() {
			if err := pin.DigitalWrite(level); err != nil {
				fmt.Println(err)
			}
			if level == 1 {
				level = 0
			} else {
				level = 1
			}
		})
	}

	robot := gobot.NewRobot("pinBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{pin},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
