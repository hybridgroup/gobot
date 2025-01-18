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
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 2(+5V), 6, 9, 14, 20 (GND)
// PWM Tinkerboard: header pin 33 (PWM2) or pin 32 (PWM3)
// Please decouple the PWM pin with an amplifier, e.g. a simple MOSFET (IRLZ34N, IRF530 etc.), a half H-bridge or
// just use a LED for the first tests.
func main() {
	adaptor := tinkerboard.NewAdaptor()
	motor := gpio.NewMotorDriver(adaptor, "32", gpio.WithMotorAnalog()) // gpio.WithMotorAnalog() is optional here

	work := func() {
		fmt.Println("first try full speed for 5 seconds...")
		if err := motor.On(); err != nil {
			log.Println(err)
		}

		time.Sleep(5 * time.Second)

		fmt.Println("second switch off for 5 seconds...")
		if err := motor.Off(); err != nil {
			log.Println(err)
		}

		time.Sleep(5 * time.Second)

		fmt.Println("finally fade in and out forever...")
		speed := byte(0)
		fadeAmount := byte(15)

		gobot.Every(100*time.Millisecond, func() {
			if err := motor.SetSpeed(speed); err != nil {
				log.Println(err)
			}
			speed = speed + fadeAmount
			if speed == 0 || speed == 255 {
				fadeAmount = -fadeAmount
			}
		})
	}

	robot := gobot.NewRobot("motorBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{motor},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
