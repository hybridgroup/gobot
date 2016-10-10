package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/leap"
)

func main() {
	gbot := gobot.NewGobot()

	leapMotionAdaptor := leap.NewAdaptor("127.0.0.1:6437")
	l := leap.NewDriver(leapMotionAdaptor)

	work := func() {
		l.On(leap.HandEvent, func(data interface{}) {
			printHand(data.(leap.Hand))
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

func printHand(hand leap.Hand) {
	fmt.Println("Hand", hand)
}
