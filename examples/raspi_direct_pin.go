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
	"gobot.io/x/gobot/v2/platforms/raspi"
)

// Wiring
// PWR  Raspi: 1 (+3.3V, VCC), 2(+5V), 6, 9, 14, 20 (GND)
// GPIO Raspi: header pin 21 (GPIO9) is input, pin 24 (GPIO8) is normal output, pin 26 (GPIO7) is inverted output
// Button: the input pin is wired with a button to GND, the internal pull up resistor is used
// LED's: the output pins are wired to the cathode of a LED, the anode is wired with a resistor (70-130Ohm for 20mA)
// to VCC
// Expected behavior: always one LED is on, the other in opposite state, if button is pressed for >2 seconds the state
// changes
func main() {
	const (
		inPinNum          = "21"
		outPinNum         = "24"
		outPinInvertedNum = "26"
	)
	// note: WithGpiosOpenDrain() is optional, if using WithGpiosOpenSource() the LED's will not light up
	board := raspi.NewAdaptor(adaptors.WithGpiosActiveLow(outPinInvertedNum),
		adaptors.WithGpiosOpenDrain(outPinNum, outPinInvertedNum),
		adaptors.WithGpiosPullUp(inPinNum),
		adaptors.WithGpioDebounce(inPinNum, 2*time.Second))

	inPin := gpio.NewDirectPinDriver(board, inPinNum)
	outPin := gpio.NewDirectPinDriver(board, outPinNum)
	outPinInverted := gpio.NewDirectPinDriver(board, outPinInvertedNum)

	work := func() {
		level := byte(1)

		gobot.Every(500*time.Millisecond, func() {
			read, err := inPin.DigitalRead()
			fmt.Printf("pin %s state is %d\n", inPinNum, read)
			if err != nil {
				fmt.Println(err)
			} else {
				level = byte(read)
			}

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

			if level == 1 {
				level = 0
			} else {
				level = 1
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
