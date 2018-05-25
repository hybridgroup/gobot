// +build example
//
// Do not build by default.

/*
How to run:
Need 2 Tello drones and 2 WiFi adapters.
Connect to the drone's Wi-Fi network from your computer. It will be named something like "TELLO-XXXXXX".

Here is the trick:
Manually setup IP address to 192.168.10.2  for first WiFi adapter, and 192.168.10.3 for second WiFi adapter.
Once you are connected to both drones, you can run the Gobot code on your computer to control the drones.

	go run examples/tello_multipledrones.go
*/

package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("192.168.10.2:8888")
	drone2 := tello.NewDriver("192.168.10.3:8888")

	work := func() {
		drone.TakeOff()
		drone2.TakeOff()

		gobot.After(5*time.Second, func() {
			drone.Land()
			drone2.Land()
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		[]gobot.Device{drone2},
		work,
	)

	robot.Start()
}
