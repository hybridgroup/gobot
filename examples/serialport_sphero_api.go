//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/drivers/serial/sphero"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
	manager := gobot.NewManager()
	api.NewAPI(manager).Start()

	spheros := map[string]string{
		"Sphero-BPO": "/dev/rfcomm0",
	}

	for name, port := range spheros {
		spheroAdaptor := serialport.NewAdaptor(port)
		spheroDriver := sphero.NewSpheroDriver(spheroAdaptor)

		work := func() {
			spheroDriver.SetRGB(uint8(255), uint8(0), uint8(0))
		}

		robot := gobot.NewRobot(name,
			[]gobot.Connection{spheroAdaptor},
			[]gobot.Device{spheroDriver},
			work,
		)
		robot.AddCommand("turn_blue", func(params map[string]interface{}) interface{} {
			spheroDriver.SetRGB(uint8(0), uint8(0), uint8(255))
			return nil
		})

		manager.AddRobot(robot)
	}

	if err := manager.Start(); err != nil {
		panic(err)
	}
}
