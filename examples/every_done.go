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
			done := gobot.Every(500*time.Millisecond, func() {
				fmt.Println("Greetings human")
			})

			gobot.After(5*time.Second, func() {
				done <- true
				fmt.Println("We're done here")
			})
		},
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
