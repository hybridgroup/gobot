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
	"gobot.io/x/gobot/v2/platforms/friendlyelec/nanopi"
)

// Wiring
// PWR  NanoPi: 1, 17 (+3.3V, VCC); 2, 4 (+5V, VDD); 6, 9, 14, 20 (GND)
// GPIO NanoPi: header pin 22 is input, pin 23 is normal output, pin 24 is inverted output
// Button: the input pin is wired with a button to GND, the internal pull up resistor is used
// LED's: the output pins are wired to the cathode of the LED, the anode is wired with a resistor (70-130Ohm for 20mA)
// to VCC
// Expected behavior: always one LED is on, the other in opposite state, if button is pressed for >2 seconds the state
// changes
func main() {
	const (
		inPinNum          = "22" // 7, 8, 10, 11, 12, 13, 15, 16, 18, 22
		outPinNum         = "23"
		outPinInvertedNum = "24"
		debounceTime      = 2 * time.Second
	)
	// note: WithGpiosOpenDrain() is optional, if using WithGpiosOpenSource() the LED's will not light up
	board := nanopi.NewNeoAdaptor(adaptors.WithGpiosActiveLow(outPinInvertedNum),
		adaptors.WithGpiosOpenDrain(outPinNum, outPinInvertedNum),
		adaptors.WithGpiosPullUp(inPinNum),
		adaptors.WithGpioDebounce(inPinNum, debounceTime))

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
				if level == 1 {
					level = 0
				} else {
					level = 1
				}
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
