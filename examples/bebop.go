package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/bebop"
)

func main() {
	gbot := gobot.NewGobot()

	bebopAdaptor := bebop.NewBebopAdaptor("Drone")
	drone := bebop.NewBebopDriver(bebopAdaptor, "Drone")

	work := func() {
		gobot.On(drone.Event("flying"), func(data interface{}) {
			gobot.After(3*time.Second, func() {
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
	gbot.AddRobot(robot)

	gbot.Start()
}
