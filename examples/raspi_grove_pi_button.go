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
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	gp := i2c.NewGrovePiDriver(r)
	button := gpio.NewButtonDriver(gp, "D3", gpio.WithButtonPollInterval(50*time.Millisecond))
	led := gpio.NewLedDriver(gp, "D2")

	work := func() {
		_ = button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("button pressed")
			if err := led.On(); err != nil {
				fmt.Println(err)
			}
		})

		_ = button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("button released")
			if err := led.Off(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{r},
		[]gobot.Device{gp, button, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
