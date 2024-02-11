//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
	go run examples/holystone_hs200.go

*/

package main

import (
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/holystone/hs200"
)

func main() {
	drone := hs200.NewDriver("172.16.10.1:8888", "172.16.10.1:8080")

	work := func() {
		drone.TakeOff()

		gobot.After(5*time.Second, func() {
			drone.Land()
		})
	}

	robot := gobot.NewRobot("hs200",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
