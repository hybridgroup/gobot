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

 This example requires that you wire an external LED to
 pin number 0 on the Microbit, as this example is intended
 to demonstrate the Microbit IOPinDriver.

 You then run the Go program on your computer and communicate
 wirelessly with the Microbit.

 How to run
 Pass the Bluetooth name or address as first param:

	go run examples/microbit_blink.go "BBC micro:bit [yowza]"

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/microbit"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])

	ubit := microbit.NewIOPinDriver(bleAdaptor)
	led := gpio.NewLedDriver(ubit, "0")

	work := func() {
		gobot.Every(1*time.Second, func() {
			if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ubit, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
