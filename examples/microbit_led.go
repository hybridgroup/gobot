// +build example
//
// Do not build by default.

/*
 How to setup
 You must be using a BBC Microbit microcontroller that has
 been flashed with the firmware from @sandeepmistry

 More info:
 https://gobot.io/documentation/platforms/microbit/#how-to-install

 This example uses the Microbit's built-in LED matrix.
 You run the Go program on your computer and communicate
 wirelessly with the Microbit.

 How to run
 Pass the Bluetooth name or address as first param:

	go run examples/microbit_led.go "BBC micro:bit [yowza]"

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/microbit"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	ubit := microbit.NewLEDDriver(bleAdaptor)

	work := func() {
		ubit.Blank()
		gobot.After(1*time.Second, func() {
			ubit.WriteText("Hello")
		})
		gobot.After(7*time.Second, func() {
			ubit.Smile()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ubit},
		work,
	)

	robot.Start()
}
