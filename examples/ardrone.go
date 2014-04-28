package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/ardrone"
)

func main() {
	ardroneAdaptor := ardrone.NewArdroneAdaptor()
	ardroneAdaptor.Name = "Drone"

	drone := ardrone.NewArdroneDriver(ardroneAdaptor)
	drone.Name = "Drone"

	work := func() {
		drone.TakeOff()
		gobot.On(drone.Events["Flying"], func(data interface{}) {
			gobot.After("1s", func() {
				drone.Right(0.1)
			})
			gobot.After("2s", func() {
				drone.Left(0.1)
			})
			gobot.After("3s", func() {
				drone.Land()
			})
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{ardroneAdaptor},
		Devices:     []gobot.Device{drone},
		Work:        work,
	}

	robot.Start()
}
