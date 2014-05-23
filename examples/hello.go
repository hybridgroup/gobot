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
			gobot.Every(0.5*time.Second, func() { fmt.Println("Greetings human") })
		},
	)

	gbot.Robots = append(gbot.Robots, robot)
	gbot.Start()
}
