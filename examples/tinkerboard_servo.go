//go:build example
// +build example

// Do not build by default.

package main

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

// Wiring
// PWR Tinkerboard: 1 (+3.3V, VCC), 2(+5V), 6, 9, 14, 20 (GND)
// PWM Tinkerboard: header pin 33 (PWM2) or pin 32 (PWM3)
func main() {
	const (
		pwmPin = "32"
		wait   = 3 * time.Second

		fiftyHzNanos = 20 * 1000 * 1000 // 50Hz = 0.02 sec = 20 ms
	)
	// usually a frequency of 50Hz is used for servos, most servos have 0.5 ms..2.5 ms for 0-180°,
	// however the mapping can be changed with options:
	adaptor := tinkerboard.NewAdaptor(
		adaptors.WithPWMDefaultPeriodForPin(pwmPin, fiftyHzNanos),
		adaptors.WithPWMServoDutyCycleRangeForPin(pwmPin, time.Millisecond, 2*time.Millisecond),
		adaptors.WithPWMServoAngleRangeForPin(pwmPin, 0, 270),
	)
	servo := gpio.NewServoDriver(adaptor, pwmPin)

	work := func() {
		fmt.Printf("first move to minimal position for %s...\n", wait)
		if err := servo.ToMin(); err != nil {
			log.Println(err)
		}

		time.Sleep(wait)

		fmt.Printf("second move to center position for %s...\n", wait)
		if err := servo.ToCenter(); err != nil {
			log.Println(err)
		}

		time.Sleep(wait)

		fmt.Printf("third move to maximal position for %s...\n", wait)
		if err := servo.ToMax(); err != nil {
			log.Println(err)
		}

		time.Sleep(wait)

		fmt.Println("finally move 0-180° (or what your servo do for the new mapping) and back forever...")
		angle := 0
		fadeAmount := 45

		gobot.Every(time.Second, func() {
			if err := servo.Move(byte(angle)); err != nil {
				log.Println(err)
			}
			angle = angle + fadeAmount
			if angle < 0 || angle > 180 {
				if angle < 0 {
					angle = 0
				}
				if angle > 180 {
					angle = 180
				}
				// change direction and recalculate
				fadeAmount = -fadeAmount
				angle = angle + fadeAmount
			}
		})
	}

	robot := gobot.NewRobot("motorBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{servo},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
