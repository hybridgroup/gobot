package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/leap"
)

func main() {
	leapMotionAdaptor := leap.NewLeapMotionAdaptor()
	leapMotionAdaptor.Name = "leap"
	leapMotionAdaptor.Port = "127.0.0.1:6437"

	leapMotionDriver := leap.NewLeapMotionDriver(leapMotionAdaptor)
	leap.Name = "leap"

	work := func() {
		gobot.On(leap.Events["Message"], func(data interface{}) {
			printHands(data.(leap.Frame))
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{leapMotionAdaptor},
		Devices:     []gobot.Device{leapMotionDriver},
		Work:        work,
	}

	robot.Start()
}

func printHands(frame leap.Frame) {
	for key, hand := range frame.Hands {
		fmt.Println("Hand", key, hand)
	}
}
