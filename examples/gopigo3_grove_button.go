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
	led := gpio.NewGroveLedDriver(gpg3, "AD_1_1")
	button := gpio.NewGroveButtonDriver(gpg3, "AD_2_1")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("On!")
			led.On()
		})
		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("Off!")
			led.Off()
		})
	}

	robot := gobot.NewRobot("gopigo3button",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{gpg3, button, led},
		work,
	)

	robot.Start()
}
