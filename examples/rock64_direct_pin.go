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
	"gobot.io/x/gobot/v2/platforms/pine64/rock64"
)

// Wiring
// PWR  ROCK64: 1, P5_1 (+3.3V, VCC); 2, 4, P5_2 (+5V, VDD); 6, 9, 14, 20, P5_7, P5_8, P5_15, P5_16 (GND)
// GPIO ROCK64: second header P5+BUS pin 3 is input, pin 4 is normal output, pin 5 is inverted output
// Button: the input pin is wired with a button to GND, the internal pull up resistor is used
// LED's: the output pins are wired to the cathode of the LED, the anode is wired with a resistor (70-130Ohm for 20mA)
// to VCC
// Expected behavior: always one LED is on, the other in opposite state, if button is pressed for >2 seconds the state
// changes
func main() {
	const (
		inPinNum          = "P5_3"
		outPinNum         = "P5_4"
		outPinInvertedNum = "P5_5"
		debounceTime      = 2 * time.Second
	)
	// note: WithGpiosOpenDrain() is optional, if using WithGpiosOpenSource() the LED's will not light up
	board := rock64.NewAdaptor(adaptors.WithGpiosActiveLow(outPinInvertedNum),
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
