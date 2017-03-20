// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/platforms/sphero"
)

func NewSwarmBot(port string) *gobot.Robot {
	spheroAdaptor := sphero.NewAdaptor(port)
	spheroDriver := sphero.NewSpheroDriver(spheroAdaptor)
	spheroDriver.SetName("Sphero" + port)

	work := func() {
		spheroDriver.Stop()

		spheroDriver.On(sphero.Collision, func(data interface{}) {
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

	master.Start()
}
