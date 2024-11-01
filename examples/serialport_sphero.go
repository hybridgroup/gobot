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
	"gobot.io/x/gobot/v2/drivers/common/spherocommon"
	"gobot.io/x/gobot/v2/drivers/serial/sphero"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
	adaptor := serialport.NewAdaptor("/dev/rfcomm0")
	spheroDriver := sphero.NewSpheroDriver(adaptor)

	work := func() {
		spheroDriver.SetDataStreaming(spherocommon.DefaultDataStreamingConfig())

		_ = spheroDriver.On(spherocommon.CollisionEvent, func(data interface{}) {
			fmt.Printf("Collision! %+v\n", data)
		})

		_ = spheroDriver.On(spherocommon.SensorDataEvent, func(data interface{}) {
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

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
