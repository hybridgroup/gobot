//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/beagleboard/beaglebone"
)

func main() {
	beagleboneAdaptor := beaglebone.NewAdaptor()
	led := gpio.NewDirectPinDriver(beagleboneAdaptor, "P8_10")
	button := gpio.NewDirectPinDriver(beagleboneAdaptor, "P8_09")

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			val, _ := button.DigitalRead()
			if val == 1 {
				if err := led.DigitalWrite(1); err != nil {
					fmt.Println(err)
				}
			} else {
				if err := led.DigitalWrite(0); err != nil {
					fmt.Println(err)
				}
			}
		})
	}

	robot := gobot.NewRobot("pinBot",
		[]gobot.Connection{beagleboneAdaptor},
		[]gobot.Device{led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
