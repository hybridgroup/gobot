package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

var Master gobot.Gobot

func TurnBlue(params map[string]interface{}) bool {
	spheroDriver := Master.FindRobotDevice(params["robotname"].(string), "sphero")
	gobot.Call(sphero.Driver, "SetRGB", uint8(0), uint8(0), uint8(255))
	return true
}

func main() {
	Master = gobot.NewGobot()
	api.Api(Master).Start()

	spheros := map[string]string{
		"Sphero-BPO": "/dev/rfcomm0",
	}

	for name, port := range spheros {
		spheroAdaptor := sphero.NewSpheroAdaptor("sphero", port)

		spheroDriver := sphero.NewSpheroDriver(spheroAdaptor, "sphero")

		work := func() {
			spheroDriver.SetRGB(uint8(255), uint8(0), uint8(0))
		}

		robot := gobot.NewRobot(name, []gobot.Connection{spheroAdaptor}, []gobot.Device{spheroDriver}, work)
		robot.Commands = map[string]interface{}{"TurnBlue": TurnBlue}

		Master.Robots = append(Master.Robots, robot)
	}

	Master.Start()
}
