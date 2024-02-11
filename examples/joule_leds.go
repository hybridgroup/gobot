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
	"gobot.io/x/gobot/v2/platforms/intel-iot/joule"
)

func main() {
	e := joule.NewAdaptor()
	led0 := gpio.NewLedDriver(e, "GP100")
	led1 := gpio.NewLedDriver(e, "GP101")
	led2 := gpio.NewLedDriver(e, "GP102")
	led3 := gpio.NewLedDriver(e, "GP103")

	work := func() {
		if err := led0.Off(); err != nil {
			fmt.Println(err)
		}
		if err := led1.Off(); err != nil {
			fmt.Println(err)
		}
		if err := led2.Off(); err != nil {
			fmt.Println(err)
		}
		if err := led3.Off(); err != nil {
			fmt.Println(err)
		}

		gobot.Every(1*time.Second, func() {
			if err := led0.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
		gobot.Every(2*time.Second, func() {
			if err := led1.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
		gobot.Every(4*time.Second, func() {
			if err := led2.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
		gobot.Every(8*time.Second, func() {
			if err := led3.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{e},
		[]gobot.Device{led0, led1, led2, led3},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
