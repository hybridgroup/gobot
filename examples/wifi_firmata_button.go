// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewTCPAdaptor(os.Args[1])
	firmataAdaptor.BoardType = "esp8266"

	button := gpio.NewButtonDriver(firmataAdaptor, "D5")
	led := gpio.NewLedDriver(firmataAdaptor, "D4")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("push")
			led.On()
		})
		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("rel")
			led.Off()
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{button, led},
		work,
	)

	robot.Start()
}
