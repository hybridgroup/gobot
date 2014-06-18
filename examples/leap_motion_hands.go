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
		gobot.On(l.Events["Message"], func(data interface{}) {
			printHands(data.(leap.Frame))
		})
	}

	gbot.Robots = append(gbot.Robots, gobot.NewRobot(
		"leapBot", []gobot.Connection{leapMotionAdaptor}, []gobot.Device{l}, work))

	gbot.Start()
}

func printHands(frame leap.Frame) {
	for key, hand := range frame.Hands {
		fmt.Println("Hand", key, hand)
	}
}
