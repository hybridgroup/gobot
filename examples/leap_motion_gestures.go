package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/leap"
)

func main() {
	gbot := gobot.NewGobot()

	leapMotionAdaptor := leap.NewLeapMotionAdaptor("leap", "127.0.0.1:6437")
	l := leap.NewLeapMotionDriver(leapMotionAdaptor, "leap")

	work := func() {
		gobot.On(l.Event("gesture"), func(data interface{}) {
			printGesture(data.(leap.Gesture))
		})
	}

	robot := gobot.NewRobot("leapBot",
		[]gobot.Connection{leapMotionAdaptor},
		[]gobot.Device{l},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}

func printGesture(gesture leap.Gesture) {
	fmt.Println("Gesture", gesture)
}
