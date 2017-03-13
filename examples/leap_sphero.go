// +build example
//
// Do not build by default.

package main

import (
	"math"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/leap"
	"gobot.io/x/gobot/platforms/sphero"
)

func main() {
	leapAdaptor := leap.NewAdaptor("127.0.0.1:6437")
	spheroAdaptor := sphero.NewAdaptor("/dev/tty.Sphero-YBW-RN-SPP")

	leapDriver := leap.NewDriver(leapAdaptor)
	spheroDriver := sphero.NewSpheroDriver(spheroAdaptor)

	work := func() {
		leapDriver.On(leap.MessageEvent, func(data interface{}) {
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

	robot.Start()
}

func scale(position float64) uint8 {
	return uint8(gobot.ToScale(gobot.FromScale(position, 0, 1), 0, 255))
}
