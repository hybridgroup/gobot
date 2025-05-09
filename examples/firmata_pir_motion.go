//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_pir_motion.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])

	sensor := gpio.NewPIRMotionDriver(firmataAdaptor, "5")
	led := gpio.NewLedDriver(firmataAdaptor, "13")

	work := func() {
		_ = sensor.On(gpio.MotionDetected, func(data interface{}) {
			fmt.Println(gpio.MotionDetected)
			if err := led.On(); err != nil {
				fmt.Println(err)
			}
		})
		_ = sensor.On(gpio.MotionStopped, func(data interface{}) {
			fmt.Println(gpio.MotionStopped)
			if err := led.Off(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("motionBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{sensor, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
