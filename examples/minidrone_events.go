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

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/parrot/minidrone"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	drone := minidrone.NewDriver(bleAdaptor)

	work := func() {
		drone.On(minidrone.Battery, func(data interface{}) {
			fmt.Printf("battery: %d\n", data)
		})

		drone.On(minidrone.FlightStatus, func(data interface{}) {
			fmt.Printf("flight status: %d\n", data)
		})

		drone.On(minidrone.Takeoff, func(data interface{}) {
			fmt.Println("taking off...")
		})

		drone.On(minidrone.Hovering, func(data interface{}) {
			fmt.Println("hovering!")
			gobot.After(5*time.Second, func() {
				drone.Land()
			})
		})

		drone.On(minidrone.Landing, func(data interface{}) {
			fmt.Println("landing...")
		})

		drone.On(minidrone.Landed, func(data interface{}) {
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
