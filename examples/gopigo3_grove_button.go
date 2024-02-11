//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/dexter/gopigo3"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	raspiAdaptor := raspi.NewAdaptor()
	gpg3 := gopigo3.NewDriver(raspiAdaptor)
	led := gpio.NewLedDriver(gpg3, "AD_1_1")
	button := gpio.NewButtonDriver(gpg3, "AD_2_1")

	work := func() {
		_ = button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("On!")
			if err := led.On(); err != nil {
				fmt.Println(err)
			}
		})
		_ = button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("Off!")
			if err := led.Off(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("gopigo3button",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{gpg3, button, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
