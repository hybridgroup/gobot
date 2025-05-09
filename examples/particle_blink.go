//go:build example
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
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/particle"
)

func main() {
	core := particle.NewAdaptor(os.Args[1], os.Args[2])
	led := gpio.NewLedDriver(core, "D7")

	work := func() {
		gobot.Every(1*time.Second, func() {
			if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		[]gobot.Device{led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
