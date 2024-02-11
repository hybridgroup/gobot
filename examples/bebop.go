//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/parrot/bebop"
)

func main() {
	bebopAdaptor := bebop.NewAdaptor()
	drone := bebop.NewDriver(bebopAdaptor)

	work := func() {
		_ = drone.On(bebop.Flying, func(data interface{}) {
			gobot.After(10*time.Second, func() {
				if err := drone.Land(); err != nil {
					fmt.Println(err)
				}
			})
		})

		if err := drone.HullProtection(true); err != nil {
			fmt.Println(err)
		}
		drone.TakeOff()
	}

	robot := gobot.NewRobot("drone",
		[]gobot.Connection{bebopAdaptor},
		[]gobot.Device{drone},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
