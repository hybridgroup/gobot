package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ardrone"
	"time"
)

func main() {
	gbot := gobot.NewGobot()
	ardroneAdaptor := ardrone.NewArdroneAdaptor("Drone")
	drone := ardrone.NewArdroneDriver(ardroneAdaptor, "Drone")

	work := func() {
		drone.TakeOff()
		gobot.On(drone.Events["Flying"], func(data interface{}) {
			gobot.After(1*time.Second, func() {
				drone.Right(0.1)
			})
			gobot.After(2*time.Second, func() {
				drone.Left(0.1)
			})
			gobot.After(3*time.Second, func() {
				drone.Land()
			})
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("drone", []gobot.Connection{ardroneAdaptor}, []gobot.Device{drone}, work))

	gbot.Start()
}
