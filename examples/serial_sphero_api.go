//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/drivers/serial"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
	master := gobot.NewMaster()
	api.NewAPI(master).Start()

	spheros := map[string]string{
		"Sphero-BPO": "/dev/rfcomm0",
	}

	for name, port := range spheros {
		spheroAdaptor := serialport.NewAdaptor(port)
		spheroDriver := serial.NewSpheroDriver(spheroAdaptor)

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

		master.AddRobot(robot)
	}

	master.Start()
}
