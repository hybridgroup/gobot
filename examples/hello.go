package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
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
