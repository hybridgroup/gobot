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
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/particle"
)

func main() {
	core := particle.NewAdaptor(os.Args[1], os.Args[2])
	led := gpio.NewLedDriver(core, "D7")
	button := gpio.NewButtonDriver(core, "D5")

	work := func() {
		button.On(button.Event("push"), func(data interface{}) {
			led.On()
		})

		button.On(button.Event("release"), func(data interface{}) {
			led.Off()
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		[]gobot.Device{button, led},
		work,
	)

	robot.Start()
}
