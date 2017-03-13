// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/leap"
)

func main() {
	leapMotionAdaptor := leap.NewAdaptor("127.0.0.1:6437")
	l := leap.NewDriver(leapMotionAdaptor)

	work := func() {
		l.On(leap.MessageEvent, func(data interface{}) {
			fmt.Println(data.(leap.Frame))
		})
	}

	robot := gobot.NewRobot("leapBot",
		[]gobot.Connection{leapMotionAdaptor},
		[]gobot.Device{l},
		work,
	)

	robot.Start()
}
