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

 This example requires that you wire an external button to
 pin number 0 on the Microbit, and also wire an external LED to
 pin number 1 on the Microbit. This example is intended
 to demonstrate using Gobot GPIO drivers with the Microbit.

 You run the Go program on your computer and communicate
 wirelessly with the Microbit.

 How to run
 Pass the Bluetooth name or address as first param:

	go run examples/microbit_io_button.go "BBC micro:bit [yowza]"

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/microbit"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])

	ubit := microbit.NewIOPinDriver(bleAdaptor)
	button := gpio.NewButtonDriver(ubit, "0")
	led := gpio.NewLedDriver(ubit, "1")

	work := func() {
		_ = button.On(gpio.ButtonPush, func(data interface{}) {
			if err := led.On(); err != nil {
				fmt.Println(err)
			}
		})
		_ = button.On(gpio.ButtonRelease, func(data interface{}) {
			if err := led.Off(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ubit, button, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
