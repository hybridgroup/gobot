// +build example
//
// Do not build by default.

/*
 How to run
 Pass the UDP port to use for the ground station to receive responses from the drone as first param:

	go run examples/tello.go "8888"
*/

package main

import (
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver(os.Args[1])

	work := func() {
		// battery, _ := drone.Battery()
		// fmt.Println("battery:", battery)

		// speed, _ := drone.Speed()
		// fmt.Println("speed:", speed)

		//		fmt.Println("Flying")
		// drone.TakeOff()

		// gobot.After(5*time.Second, func() {
		// 	drone.Land()
		// 	ft, _ := drone.FlightTime()
		// 	fmt.Println("flight time:", ft)
		// })
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
