// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
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

	robot.Start()
}
