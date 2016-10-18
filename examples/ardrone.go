package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ardrone"
)

func main() {
	ardroneAdaptor := ardrone.NewAdaptor()
	drone := ardrone.NewDriver(ardroneAdaptor)

	work := func() {
		drone.On(ardrone.Flying, func(data interface{}) {
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

	robot.Start()
}
