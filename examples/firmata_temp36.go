//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_temp36.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])

	work := func() {
		gobot.Every(1*time.Second, func() {
			val, err := firmataAdaptor.AnalogRead("0")
			if err != nil {
				fmt.Println(err)
				return
			}

			voltage := (float64(val) * 5) / 1024 // if using 3.3V replace 5 with 3.3
			tempC := (voltage - 0.5) * 100
			tempF := (tempC * 9 / 5) + 32

			fmt.Printf("%.2f°C\n", tempC)
			fmt.Printf("%.2f°F\n", tempF)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
