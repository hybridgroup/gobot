// +build example
//
// Do not build by default.

/*
 How to run
 Pass the IP address for the ground station as first param:

	go run examples/tello.go "192.168.10.2:8888"
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver(os.Args[1])

	work := func() {
		fmt.Println("Flying")
		drone.TakeOff()

		gobot.After(5*time.Second, func() {
			drone.Land()
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
