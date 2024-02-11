//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/chip"
)

func main() {
	chipAdaptor := chip.NewAdaptor()
	button := gpio.NewButtonDriver(chipAdaptor, "XIO-P6")
	led := gpio.NewLedDriver(chipAdaptor, "XIO-P7")

	work := func() {
		_ = button.On(gpio.ButtonPush, func(data interface{}) {
			if err := led.On(); err != nil {
				fmt.Println(err)
			}
		})

		_ = button.On(gpio.ButtonRelease, func(data interface{}) {
			if err := led.Off(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{chipAdaptor},
		[]gobot.Device{button, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
