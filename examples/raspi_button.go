// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	button := gpio.NewButtonDriver(r, "11")
	led := gpio.NewLedDriver(r, "7")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("button pressed")
			led.On()
		})

		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("button released")
			led.Off()
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{r},
		[]gobot.Device{button, led},
		work,
	)

	robot.Start()
}
