// +build example
//
// Do not build by default.

/*
 To run this example, pass the device ID as first param,
 and the access token as the second param:

	go run examples/particle_blink.go mydevice myaccesstoken
*/

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/particle"
)

func main() {
	core := particle.NewAdaptor(os.Args[1], os.Args[2])
	led := gpio.NewLedDriver(core, "D7")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
