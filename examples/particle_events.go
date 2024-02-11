//go:build example
// +build example

//
// Do not build by default.

/*
 To run this example, pass the device ID as first param,
 and the access token as the second param:

	go run examples/particle_events.go mydevice myaccesstoken
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/particle"
)

func main() {
	core := particle.NewAdaptor(os.Args[1], os.Args[2])

	work := func() {
		if stream, err := core.EventStream("all", ""); err != nil {
			fmt.Println(err)
		} else {
			// TODO: some other way to handle this
			fmt.Println(stream)
		}
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
