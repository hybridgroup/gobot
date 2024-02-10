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
	"gobot.io/x/gobot/v2/platforms/upboard/up2"
)

func main() {
	b := up2.NewAdaptor()
	red := gpio.NewLedDriver(b, up2.LEDRed)
	blue := gpio.NewLedDriver(b, up2.LEDBlue)
	green := gpio.NewLedDriver(b, up2.LEDGreen)
	yellow := gpio.NewLedDriver(b, up2.LEDYellow)

	work := func() {
		if err := red.Off(); err != nil {
			fmt.Println(err)
		}
		if err := blue.Off(); err != nil {
			fmt.Println(err)
		}
		if err := green.Off(); err != nil {
			fmt.Println(err)
		}
		if err := yellow.Off(); err != nil {
			fmt.Println(err)
		}

		gobot.Every(1*time.Second, func() {
			if err := red.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
		gobot.Every(2*time.Second, func() {
			if err := green.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
		gobot.Every(4*time.Second, func() {
			if err := yellow.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
		gobot.Every(8*time.Second, func() {
			if err := blue.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{b},
		[]gobot.Device{red, blue, green, yellow},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
