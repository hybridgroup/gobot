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
		l.On(leap.GestureEvent, func(data interface{}) {
			printGesture(data.(leap.Gesture))
		})
	}

	robot := gobot.NewRobot("leapBot",
		[]gobot.Connection{leapMotionAdaptor},
		[]gobot.Device{l},
		work,
	)

	robot.Start()
}

func printGesture(gesture leap.Gesture) {
	fmt.Println("Gesture", gesture)
}
