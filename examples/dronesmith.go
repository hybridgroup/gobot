// to run:
// go run examples/dronesmith.go droneid email key
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/dronesmith"
)

func main() {
	a := dronesmith.NewAdaptor(os.Args[1], os.Args[2], os.Args[3])
	telemetry := dronesmith.NewTelemetryDriver(a)
	control := dronesmith.NewControlDriver(a)

	work := func() {
		fmt.Println(telemetry.Info())
		fmt.Println("arming...")
		control.Arm()
		fmt.Println("taking off...")
		control.Takeoff()
		gobot.After(10*time.Second, func() {
			fmt.Println("landing...")
			control.Land()
		})
	}

	robot := gobot.NewRobot("mydrone",
		[]gobot.Connection{a},
		[]gobot.Device{telemetry, control},
		work,
	)

	robot.Start()
}
