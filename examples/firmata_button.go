//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_button.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])

	button := gpio.NewButtonDriver(firmataAdaptor, "2")
	led := gpio.NewLedDriver(firmataAdaptor, "3")

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
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{button, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
