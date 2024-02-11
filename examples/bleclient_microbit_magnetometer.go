//go:build example
// +build example

//
// Do not build by default.

/*
 How to setup
 You must be using a BBC Microbit microcontroller that has
 been flashed with the firmware from @sandeepmistry

 More info:
 https://gobot.io/documentation/platforms/microbit/#how-to-install

 This example uses the Microbit's built-in magnetometer.
 You run the Go program on your computer and communicate
 wirelessly with the Microbit.

 How to run
 Pass the Bluetooth name or address as first param:

	go run examples/microbit_magnetometer.go "BBC micro:bit [yowza]"

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/microbit"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])
	ubit := microbit.NewMagnetometerDriver(bleAdaptor)

	work := func() {
		_ = ubit.On(microbit.MagnetometerEvent, func(data interface{}) {
			fmt.Println("Magnetometer", data)
		})
	}

	robot := gobot.NewRobot("magnetoBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ubit},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
