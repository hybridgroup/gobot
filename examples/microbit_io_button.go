// +build example
//
// Do not build by default.

package main

import (
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/microbit"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])

	ubit := microbit.NewIOPinDriver(bleAdaptor)
	button := gpio.NewButtonDriver(ubit, "0")
	led := gpio.NewLedDriver(ubit, "1")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			led.On()
		})
		button.On(gpio.ButtonRelease, func(data interface{}) {
			led.Off()
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{bleAdaptor, ubit},
		[]gobot.Device{ubit, button, led},
		work,
	)

	robot.Start()
}
