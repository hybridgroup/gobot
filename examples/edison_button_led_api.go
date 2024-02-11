//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/intel-iot/edison"
)

func main() {
	master := gobot.NewMaster()
	api.NewAPI(master).Start()

	e := edison.NewAdaptor()

	button := gpio.NewButtonDriver(e, "2")
	led := gpio.NewLedDriver(e, "4")

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
		[]gobot.Connection{e},
		[]gobot.Device{led, button},
		work,
	)

	master.AddRobot(robot)

	if err := master.Start(); err != nil {
		panic(err)
	}
}
