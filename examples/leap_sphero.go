//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"math"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/serial/sphero"
	"gobot.io/x/gobot/v2/platforms/leap"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
	leapAdaptor := leap.NewAdaptor("127.0.0.1:6437")
	spheroAdaptor := serialport.NewAdaptor("/dev/tty.Sphero-YBW-RN-SPP")

	leapDriver := leap.NewDriver(leapAdaptor)
	spheroDriver := sphero.NewSpheroDriver(spheroAdaptor)

	work := func() {
		_ = leapDriver.On(leap.MessageEvent, func(data interface{}) {
			hands := data.(leap.Frame).Hands

			if len(hands) > 0 {
				x := math.Abs(hands[0].Direction[0])
				y := math.Abs(hands[0].Direction[1])
				z := math.Abs(hands[0].Direction[2])
				spheroDriver.SetRGB(scale(x), scale(y), scale(z))
			}
		})
	}

	robot := gobot.NewRobot("leapBot",
		[]gobot.Connection{leapAdaptor, spheroAdaptor},
		[]gobot.Device{leapDriver, spheroDriver},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}

func scale(position float64) uint8 {
	return uint8(gobot.ToScale(gobot.FromScale(position, 0, 1), 0, 255))
}
