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
		nil,
		nil,
		func() {
			gobot.Every(500*time.Millisecond, func() { fmt.Println("Greetings human") })
		},
	)

	gbot.Robots = append(gbot.Robots, robot)
	gbot.Start()
}
