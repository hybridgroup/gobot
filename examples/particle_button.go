//go:build example
// +build example

//
// Do not build by default.

/*
 To run this example, pass the device ID as first param,
 and the access token as the second param:

	go run examples/particle_button.go mydevice myaccesstoken
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/particle"
)

func main() {
	core := particle.NewAdaptor(os.Args[1], os.Args[2])
	led := gpio.NewLedDriver(core, "D7")
	button := gpio.NewButtonDriver(core, "D5")

	work := func() {
		_ = button.On(button.Event("push"), func(data interface{}) {
			if err := led.On(); err != nil {
				fmt.Println(err)
			}
		})

		_ = button.On(button.Event("release"), func(data interface{}) {
			if err := led.Off(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		[]gobot.Device{button, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
