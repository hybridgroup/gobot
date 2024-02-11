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

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/microbit"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])
	ubit := microbit.NewButtonDriver(bleAdaptor)

	work := func() {
		_ = ubit.On(microbit.ButtonAEvent, func(data interface{}) {
			fmt.Println("button A", data)
		})

		_ = ubit.On(microbit.ButtonBEvent, func(data interface{}) {
			fmt.Println("button B", data)
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ubit},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
