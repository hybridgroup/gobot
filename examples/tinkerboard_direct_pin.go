// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/platforms/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 2(+5V), 6, 9, 14, 20 (GND)
// GPIO Tinkerboard: header pin 26 used as normal output, pin 27 used as inverted output
func main() {
	const (
		pinNum         = "7"
		pinInvertedNum = "22"
	)
	board := tinkerboard.NewAdaptor(adaptors.WithGpiosActiveLow(pinInvertedNum))
	pin := gpio.NewDirectPinDriver(board, pinNum)
	pinInverted := gpio.NewDirectPinDriver(board, pinInvertedNum)

	work := func() {
		level := byte(1)

		gobot.Every(500*time.Millisecond, func() {
			err := pin.DigitalWrite(level)
			fmt.Printf("pin %s is now %d\n", pinNum, level)
			if err != nil {
				fmt.Println(err)
			}

			err = pinInverted.DigitalWrite(level)
			fmt.Printf("pin %s is now not %d\n", pinInvertedNum, level)
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
