//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/drivers/common/sphero"
	"gobot.io/x/gobot/v2/drivers/serial"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func NewSwarmBot(port string) *gobot.Robot {
	spheroAdaptor := serialport.NewAdaptor(port)
	spheroDriver := serial.NewSpheroDriver(spheroAdaptor, serial.WithName("Sphero"+port))

	work := func() {
		spheroDriver.Stop()

		_ = spheroDriver.On(sphero.CollisionEvent, func(data interface{}) {
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
	master := gobot.NewMaster()
	api.NewAPI(master).Start()

	spheros := []string{
		"/dev/rfcomm0",
		"/dev/rfcomm1",
		"/dev/rfcomm2",
		"/dev/rfcomm3",
	}

	for _, port := range spheros {
		master.AddRobot(NewSwarmBot(port))
	}

	if err := master.Start(); err != nil {
		panic(err)
	}
}
