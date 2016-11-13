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

	work := func() {
		fmt.Println(a.Info())
	}

	robot := gobot.NewRobot("mydrone",
		[]gobot.Connection{a},
		work,
	)

	robot.Start()
}
