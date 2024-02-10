//go:build example
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
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewTCPAdaptor(os.Args[1])
	led := gpio.NewLedDriver(firmataAdaptor, "2")

	work := func() {
		gobot.Every(1*time.Second, func() {
			if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
