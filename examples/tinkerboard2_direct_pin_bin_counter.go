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
	"gobot.io/x/gobot/v2/platforms/tinkerboard/tinkerboard2"
)

// Wiring
// PWR  Tinkerboard-2: 1 (+3.3V, VCC), 2(+5V), 6, 9, 14, 20 (GND)
// GPIO Tinkerboard-2: header pins 3, 5, 7, 11 used as inverted output
// LED's: the output pins are wired to the cathode of a LED, the anode is wired with a resistor (70-130Ohm for 20mA)
// to +3.3V (use >150Ohm if connected to +5V)
// Expected behavior: the 4 LED's on normal output counts up binary
func main() {
	const (
		outPinBit0Num = "3"
		outPinBit1Num = "5"
		outPinBit2Num = "7"
		outPinBit3Num = "11"
	)

	board := tinkerboard2.NewAdaptor(adaptors.WithGpiosActiveLow(outPinBit0Num, outPinBit1Num, outPinBit2Num,
		outPinBit3Num))
	outPinB0 := gpio.NewDirectPinDriver(board, outPinBit0Num)
	outPinB1 := gpio.NewDirectPinDriver(board, outPinBit1Num)
	outPinB2 := gpio.NewDirectPinDriver(board, outPinBit2Num)
	outPinB3 := gpio.NewDirectPinDriver(board, outPinBit3Num)

	work := func() {
		value := byte(0)

		gobot.Every(500*time.Millisecond, func() {
			b0 := value & 0x01
			b1 := (value & 0x02) / 0x02
			b2 := (value & 0x04) / 0x04
			b3 := (value & 0x08) / 0x08

			if err := outPinB0.DigitalWrite(b0); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("pin %s is now %d\n", outPinBit0Num, b0)
			}

			if err := outPinB1.DigitalWrite(b1); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("pin %s is now %d\n", outPinBit1Num, b1)
			}

			if err := outPinB2.DigitalWrite(b2); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("pin %s is now %d\n", outPinBit2Num, b2)
			}

			if err := outPinB3.DigitalWrite(b3); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("pin %s is now %d\n", outPinBit3Num, b3)
			}

			value++
			if value > 15 {
				value = 0
			}
		})
	}

	robot := gobot.NewRobot("pinBot",
		[]gobot.Connection{board},
		[]gobot.Device{outPinB0, outPinB1, outPinB2, outPinB3},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
