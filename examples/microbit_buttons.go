// +build example
//
// Do not build by default.

/*
 How to setup
 You must be using a BBC Microbit microcontroller that has
 been flashed with the firmware from @sandeepmistry

 More info:
 https://gobot.io/documentation/platforms/microbit/#how-to-install

 This example uses the Microbit's built-in buttons. You run
 the Go program on your computer and communicate wirelessly
 with the Microbit.

 How to run
 Pass the Bluetooth name or address as first param:

	go run examples/microbit_buttons_led.go "BBC micro:bit [yowza]"
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
	ubit := microbit.NewButtonDriver(bleAdaptor)

	work := func() {
		ubit.On(microbit.ButtonA, func(data interface{}) {
			fmt.Println("button A", data)
		})

		ubit.On(microbit.ButtonB, func(data interface{}) {
			fmt.Println("button B", data)
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ubit},
		work,
	)

	robot.Start()
}
