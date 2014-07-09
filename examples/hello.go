package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"time"
)

func main() {
	gbot := gobot.NewGobot()

	robot := gobot.NewRobot(
		"hello",
		func() {
			gobot.Every(500*time.Millisecond, func() { fmt.Println("Greetings human") })
		},
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
