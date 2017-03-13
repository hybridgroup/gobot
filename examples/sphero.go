// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/sphero"
)

func main() {
	adaptor := sphero.NewAdaptor("/dev/rfcomm0")
	spheroDriver := sphero.NewSpheroDriver(adaptor)

	work := func() {
		spheroDriver.SetDataStreaming(sphero.DefaultDataStreamingConfig())

		spheroDriver.On(sphero.Collision, func(data interface{}) {
			fmt.Printf("Collision! %+v\n", data)
		})

		spheroDriver.On(sphero.SensorData, func(data interface{}) {
			fmt.Printf("Streaming Data! %+v\n", data)
		})

		gobot.Every(3*time.Second, func() {
			spheroDriver.Roll(30, uint16(gobot.Rand(360)))
		})

		gobot.Every(1*time.Second, func() {
			r := uint8(gobot.Rand(255))
			g := uint8(gobot.Rand(255))
			b := uint8(gobot.Rand(255))
			spheroDriver.SetRGB(r, g, b)
		})
	}

	robot := gobot.NewRobot("sphero",
		[]gobot.Connection{adaptor},
		[]gobot.Device{spheroDriver},
		work,
	)

	robot.Start()
}
