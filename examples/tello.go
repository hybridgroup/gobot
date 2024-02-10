//go:build example
// +build example

//
// Do not build by default.

/*
How to run:
Connect to the drone's Wi-Fi network from your computer. It will be named something like "TELLO-XXXXXX".

Once you are connected you can run the Gobot code on your computer to control the drone.

	go run examples/tello.go
*/

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8888")

	work := func() {
		if err := drone.TakeOff(); err != nil {
			fmt.Println(err)
		}

		gobot.After(5*time.Second, func() {
			if err := drone.Land(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
