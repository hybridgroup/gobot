//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 2(+5V), 6, 9, 14, 20 (GND)
// GPIO Tinkerboard: header pin 7 is the trigger output, pin 22 used as echo input
// HC-SR04: the power is wired to +5V and GND of tinkerboard, the same for trigger output and the echo input pin
func main() {
	const (
		triggerOutput = "7"
		echoInput     = "22"
	)

	// this is mandatory for systems with defunct edge detection, although the "cdev" is used with an newer Kernel
	// keep in mind, that this cause more inaccurate measurements
	const pollEdgeDetection = true
	opts := []interface{}{}
	if pollEdgeDetection {
		opts = append(opts, gpio.WithHCSR04UseEdgePolling())
	}

	a := tinkerboard.NewAdaptor()
	hcsr04 := gpio.NewHCSR04Driver(a, triggerOutput, echoInput, opts...)

	work := func() {
		if pollEdgeDetection {
			fmt.Println("Please note that measurements are CPU consuming and will be more inaccurate with this setting.")
			fmt.Println("After startup the system is under load and the measurement is very inaccurate, so wait a bit...")
			time.Sleep(2000 * time.Millisecond)
		}

		if err := hcsr04.StartDistanceMonitor(); err != nil {
			log.Fatal(err)
		}

		// first single shot
		if v, err := hcsr04.MeasureDistance(); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("first single shot done: %5.3f m\n", v)
		}

		ticker := gobot.Every(1*time.Second, func() {
			fmt.Printf("continuous measurement: %5.3f m\n", hcsr04.Distance())
		})

		gobot.After(5*time.Second, func() {
			if err := hcsr04.StopDistanceMonitor(); err != nil {
				log.Fatal(err)
			}
			ticker.Stop()
		})

		gobot.After(7*time.Second, func() {
			// second single shot
			if v, err := hcsr04.MeasureDistance(); err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("second single shot done: %5.3f m\n", v)
			}
			// cleanup
			if err := hcsr04.Halt(); err != nil {
				log.Println(err)
			}
			if err := a.Finalize(); err != nil {
				log.Println(err)
			}
			os.Exit(0)
		})
	}

	robot := gobot.NewRobot("distanceBot",
		[]gobot.Connection{a},
		[]gobot.Device{hcsr04},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
