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
		gobot.On(drone.Event("Flying"), func(data interface{}) {
			gobot.After(3*time.Second, func() {
				drone.Land()
			})
		})
	}

	robot := gobot.NewRobot("drone",
		[]gobot.Connection{ardroneAdaptor},
		[]gobot.Device{drone},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
