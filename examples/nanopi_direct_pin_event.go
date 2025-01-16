//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/friendlyelec/nanopi"
	"gobot.io/x/gobot/v2/system"
)

const (
	inPinNum          = "22" // 7, 8, 10, 11, 12, 13, 15, 16, 18, 22
	outPinNum         = "23"
	outPinInvertedNum = "24"
	debounceTime      = 2 * time.Second
)

var (
	outPin         gobot.DigitalPinner
	outPinInverted gobot.DigitalPinner
)

// Wiring
// PWR  NanoPi: 1, 17 (+3.3V, VCC); 2, 4 (+5V, VDD); 6, 9, 14, 20 (GND)
// GPIO NanoPi: header pin 22 is input, pin 23 is normal output, pin 24 is inverted output
// Button: the input pin is wired with a button to GND, the internal pull up resistor is used
// LED's: the output pins are wired to the cathode of a LED, the anode is wired with a resistor (70-130Ohm for 20mA)
// to VCC
// Expected behavior: always one LED is on, the other in opposite state, if button is pressed for >2 seconds the state
// changes
func main() {
	board := nanopi.NewNeoAdaptor()

	work := func() {
		inPin, err := board.DigitalPin(inPinNum)
		if err != nil {
			fmt.Println(err)
		}
		if err := inPin.ApplyOptions(system.WithPinDirectionInput(), system.WithPinPullUp(),
			system.WithPinDebounce(debounceTime), system.WithPinEventOnBothEdges(buttonEventHandler)); err != nil {
			fmt.Println(err)
		}

		// note: WithPinOpenDrain() is optional, if using WithPinOpenSource() the LED's will not light up
		outPin, err = board.DigitalPin(outPinNum)
		if err != nil {
			fmt.Println(err)
		}
		if err := outPin.ApplyOptions(system.WithPinDirectionOutput(1), system.WithPinOpenDrain()); err != nil {
			fmt.Println(err)
		}

		outPinInverted, err = board.DigitalPin(outPinInvertedNum)
		if err != nil {
			fmt.Println(err)
		}
		if err := outPinInverted.ApplyOptions(system.WithPinActiveLow(), system.WithPinDirectionOutput(1),
			system.WithPinOpenDrain()); err != nil {
			fmt.Println(err)
		}

		fmt.Printf("\nPlease press and hold the button for at least %s\n", debounceTime)
	}

	robot := gobot.NewRobot("pinEdgeBot",
		[]gobot.Connection{board},
		[]gobot.Device{},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}

func buttonEventHandler(offset int, t time.Duration, et string, sn uint32, lsn uint32) {
	fmt.Printf("%s: %s detected on line %d with total sequence %d and line sequence %d\n", t, et, offset, sn, lsn)
	level := 1

	if et == "falling edge" {
		level = 0
	}

	if outPin != nil {
		err := outPin.Write(level)
		fmt.Printf("pin %s is now %d\n", outPinNum, level)
		if err != nil {
			fmt.Println(err)
		}
	}

	if outPinInverted != nil {
		err := outPinInverted.Write(level)
		fmt.Printf("pin %s is now not %d\n", outPinInvertedNum, level)
		if err != nil {
			fmt.Println(err)
		}
	}
}
