// to run:
// go run examples/dronesmith.go droneid email key
package main

import (
	"fmt"
	"os"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/dronesmith"
)

func main() {
	a := dronesmith.NewAdaptor(os.Args[1], os.Args[2], os.Args[3])
	tel := dronesmith.NewTelemetryDriver(a)

	work := func() {
		fmt.Println(tel.Info())
	}

	robot := gobot.NewRobot("mydrone",
		[]gobot.Connection{a},
		[]gobot.Device{tel},
		work,
	)

	robot.Start()
}
