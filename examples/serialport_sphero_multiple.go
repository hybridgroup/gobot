//go:build example
// +build example

//
// Do not build by default.

//nolint:gosec // ok here
package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/drivers/common/spherocommon"
	"gobot.io/x/gobot/v2/drivers/serial"
	"gobot.io/x/gobot/v2/drivers/serial/sphero"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func NewSwarmBot(port string) *gobot.Robot {
	spheroAdaptor := serialport.NewAdaptor(port)
	spheroDriver := sphero.NewSpheroDriver(spheroAdaptor, serial.WithName("Sphero"+port))

	work := func() {
		spheroDriver.Stop()

		_ = spheroDriver.On(spherocommon.CollisionEvent, func(data interface{}) {
			fmt.Println("Collision Detected!")
		})

		gobot.Every(1*time.Second, func() {
			spheroDriver.Roll(100, uint16(gobot.Rand(360)))
		})
		gobot.Every(3*time.Second, func() {
			spheroDriver.SetRGB(uint8(gobot.Rand(255)),
				uint8(gobot.Rand(255)),
				uint8(gobot.Rand(255)),
			)
		})
	}

	robot := gobot.NewRobot("sphero",
		[]gobot.Connection{spheroAdaptor},
		[]gobot.Device{spheroDriver},
		work,
	)

	return robot
}

func main() {
	manager := gobot.NewManager()
	api.NewAPI(manager).Start()

	spheros := []string{
		"/dev/rfcomm0",
		"/dev/rfcomm1",
		"/dev/rfcomm2",
		"/dev/rfcomm3",
	}

	for _, port := range spheros {
		manager.AddRobot(NewSwarmBot(port))
	}

	if err := manager.Start(); err != nil {
		panic(err)
	}
}
