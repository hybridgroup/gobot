//go:build example
// +build example

// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 2(+5V), 6, 9, 14, 20 (GND)
// GPIO Tinkerboard: header pin 22 is input, pin 7 is normal output
// Button: the input pin is wired with a button to GND, an external pull up resistor is needed (e.g. 1K)
// LED: the output pin is wired to the cathode of a LED, the anode is wired with a resistor (70-130Ohm for 20mA) to VCC
// Expected behavior: LED is initially on, if button is pressed and released, the state changes
func main() {
	const (
		buttonPin = "22"
		ledPin    = "7"
	)

	a := tinkerboard.NewAdaptor()
	button := gpio.NewButtonDriver(a, buttonPin, gpio.WithButtonPollInterval(100*time.Millisecond))
	led := gpio.NewLedDriver(a, ledPin)
	if err := led.On(); err != nil {
		fmt.Println(err)
	}

	work := func() {
		if err := button.On(gpio.Error, func(err interface{}) {
			fmt.Println("an error occurred:", err)
		}); err != nil {
			panic(err)
		}

		if err := button.On(gpio.ButtonPush, func(interface{}) {
			fmt.Println("button pressed")
			if err := led.On(); err != nil {
				fmt.Println(err)
			}
		}); err != nil {
			panic(err)
		}

		if err := button.On(gpio.ButtonRelease, func(interface{}) {
			fmt.Println("button released")
			if err := led.Off(); err != nil {
				fmt.Println(err)
			}
		}); err != nil {
			panic(err)
		}
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{a},
		[]gobot.Device{button, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
