package main

import (
	"math"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/leap"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

func main() {
	gbot := gobot.NewGobot()

	leapAdaptor := leap.NewLeapMotionAdaptor("leap", "127.0.0.1:6437")
	spheroAdaptor := sphero.NewSpheroAdaptor("Sphero", "/dev/tty.Sphero-YBW-RN-SPP")

	leapDriver := leap.NewLeapMotionDriver(leapAdaptor, "leap")
	spheroDriver := sphero.NewSpheroDriver(spheroAdaptor, "sphero")

	work := func() {
		gobot.On(leapDriver.Event("message"), func(data interface{}) {
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

	gbot.AddRobot(robot)

	gbot.Start()
}

func scale(position float64) uint8 {
	return uint8(gobot.ToScale(gobot.FromScale(position, 0, 1), 0, 255))
}
