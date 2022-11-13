// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 2(+5V), 6, 9, 14, 20 (GND)
// GPIO Tinkerboard: header pin 26 used as output
func main() {
	const pinNo = "26"
	board := tinkerboard.NewAdaptor()
	pin := gpio.NewDirectPinDriver(board, pinNo)

	work := func() {
		level := byte(1)

		gobot.Every(500*time.Millisecond, func() {
			err := pin.DigitalWrite(level)
			fmt.Printf("pin %s is now %d\n", pinNo, level)
			if err != nil {
				fmt.Println(err)
			}
			if level == 1 {
				level = 0
			} else {
				level = 1
			}

		})
	}

	robot := gobot.NewRobot("pinBot",
		[]gobot.Connection{board},
		[]gobot.Device{pin},
		work,
	)

	robot.Start()
}
