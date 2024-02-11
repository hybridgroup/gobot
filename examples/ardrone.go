//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/parrot/ardrone"
)

func main() {
	ardroneAdaptor := ardrone.NewAdaptor()
	drone := ardrone.NewDriver(ardroneAdaptor)

	work := func() {
		_ = drone.On(ardrone.Flying, func(data interface{}) {
			gobot.After(3*time.Second, func() {
				drone.Land()
			})
		})
		drone.TakeOff()
	}

	robot := gobot.NewRobot("drone",
		[]gobot.Connection{ardroneAdaptor},
		[]gobot.Device{drone},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
