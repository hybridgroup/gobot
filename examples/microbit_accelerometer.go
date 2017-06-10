// +build example
//
// Do not build by default.

/*
 How to setup
 You must be using a BBC Microbit microcontroller that has
 been flashed with the firmware from @sandeepmistry

 More info:
 https://gobot.io/documentation/platforms/microbit/#how-to-install

 This example uses the Microbit's built-in accelerometer.
 You run the Go program on your computer and communicate
 wirelessly with the Microbit.

 How to run
 Pass the Bluetooth name or address as first param:

	go run examples/microbit_accelerometer.go "BBC micro:bit [yowza]"

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/microbit"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	ubit := microbit.NewAccelerometerDriver(bleAdaptor)

	work := func() {
		ubit.On(microbit.Accelerometer, func(data interface{}) {
			fmt.Println("Accelerometer", data)
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ubit},
		work,
	)

	robot.Start()
}
