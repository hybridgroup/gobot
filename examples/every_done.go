//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
)

func main() {
	robot := gobot.NewRobot(
		"hello",
		func() {
			done := gobot.Every(750*time.Millisecond, func() {
				fmt.Println("Greetings human")
			})

			gobot.After(5*time.Second, func() {
				done.Stop()
				fmt.Println("We're done here")
			})
		},
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
