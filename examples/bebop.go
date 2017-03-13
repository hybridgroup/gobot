// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/parrot/bebop"
)

func main() {
	bebopAdaptor := bebop.NewAdaptor()
	drone := bebop.NewDriver(bebopAdaptor)

	work := func() {
		drone.On(bebop.Flying, func(data interface{}) {
			gobot.After(10*time.Second, func() {
				drone.Land()
			})
		})

		drone.HullProtection(true)
		drone.TakeOff()
	}

	robot := gobot.NewRobot("drone",
		[]gobot.Connection{bebopAdaptor},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
