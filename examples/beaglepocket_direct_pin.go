//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/platforms/beaglebone"
)

// Wiring
// PWR  Pocket: P1.14, P2.23 (+3.3V, VCC); P1.15, P1.16, P1.22, P2.15, P2.21 (GND)
// GPIO Pocket: header pin P1.34 is input, pin P1.35 is normal output, pin P1.36 is inverted output
// Button: the input pin is wired with a button to GND, an external pull up resistor is needed (e.g. 2kOhm to VCC)
// LED's: the output pins are wired to the cathode of the LED, the anode is wired with a resistor (70-130Ohm for 20mA)
// to VCC
// Expected behavior: always one LED is on, the other in opposite state, if button is pressed the state changes
func main() {
	const (
		inPinNum          = "P1_34"
		outPinNum         = "P1_35"
		outPinInvertedNum = "P1_36"
	)

	board := beaglebone.NewPocketBeagleAdaptor(adaptors.WithGpiosActiveLow(outPinInvertedNum))

	inPin := gpio.NewDirectPinDriver(board, inPinNum)
	outPin := gpio.NewDirectPinDriver(board, outPinNum)
	outPinInverted := gpio.NewDirectPinDriver(board, outPinInvertedNum)

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			read, err := inPin.DigitalRead()
			fmt.Printf("pin %s state is %d\n", inPinNum, read)
			if err != nil {
				fmt.Println(err)
			}

			level := byte(read)

			err = outPin.DigitalWrite(level)
			fmt.Printf("pin %s is now %d\n", outPinNum, level)
			if err != nil {
				fmt.Println(err)
			}

			err = outPinInverted.DigitalWrite(level)
			fmt.Printf("pin %s is now not %d\n", outPinInvertedNum, level)
			if err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("pinBot",
		[]gobot.Connection{board},
		[]gobot.Device{inPin, outPin, outPinInverted},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
