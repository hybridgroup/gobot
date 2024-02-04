//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass the Bluetooth name or address as first param:

	go run examples/minidrone_events.go "Mambo_1234"

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/parrot"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])
	drone := parrot.NewMinidroneDriver(bleAdaptor)

	work := func() {
		drone.On(parrot.BatteryEvent, func(data interface{}) {
			fmt.Printf("battery: %d\n", data)
		})

		drone.On(parrot.FlightStatusEvent, func(data interface{}) {
			fmt.Printf("flight status: %d\n", data)
		})

		drone.On(parrot.TakeoffEvent, func(data interface{}) {
			fmt.Println("taking off...")
		})

		drone.On(parrot.HoveringEvent, func(data interface{}) {
			fmt.Println("hovering!")
			gobot.After(5*time.Second, func() {
				drone.Land()
			})
		})

		drone.On(parrot.LandingEvent, func(data interface{}) {
			fmt.Println("landing...")
		})

		drone.On(parrot.LandedEvent, func(data interface{}) {
			fmt.Println("landed.")
		})

		time.Sleep(1000 * time.Millisecond)
		drone.TakeOff()
	}

	robot := gobot.NewRobot("minidrone",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
