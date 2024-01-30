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

 This example uses the Microbit's built-in buttons and
 built-in LED matrix. You run the Go program on your computer and
 communicate wirelessly with the Microbit.

 How to run
 Pass the Bluetooth name or address as first param:

	go run examples/microbit_buttons_led.go "BBC micro:bit [yowza]"

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"os"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/microbit"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])
	buttons := microbit.NewButtonDriver(bleAdaptor)
	leds := microbit.NewLEDDriver(bleAdaptor)

	work := func() {
		buttons.On(microbit.ButtonAEvent, func(data interface{}) {
			if data.([]byte)[0] == 1 {
				leds.UpLeftArrow()
				return
			}

			leds.Blank()
		})

		buttons.On(microbit.ButtonBEvent, func(data interface{}) {
			if data.([]byte)[0] == 1 {
				leds.UpRightArrow()
				return
			}

			leds.Blank()
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{buttons, leds},
		work,
	)

	robot.Start()
}
