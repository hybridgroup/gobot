// +build example
//
// Do not build by default.

/*
 How to setup
 You must be using a WiFi-connected microcontroller,
 that has been flashed with the WifiFirmata sketch.
 You then run the go program on your computer, and communicate
 wirelessly with the microcontroller.

 How to run
 Pass the IP address/port as first param:

	go run examples/wifi_firmata_blink.go 192.168.0.39:3030
*/

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewTCPAdaptor(os.Args[1])
	led := gpio.NewLedDriver(firmataAdaptor, "2")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
