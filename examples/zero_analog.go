//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/platforms/radxa/zero"
)

// Wiring:
// PWR            : 1, 17 (+3.3V, VCC), 2, 4 (+5V), 6, 9, 14, 20, 25, 30, 34, 39 (GND)
// ADC (max. 1.8V): header pin 15 is input for channel 1, pin 26 is input for channel 2
func main() {
	const (
		inPin0         = "15_mean"
		inPin1         = "26"
		inVoltageScale = 0.439453125 // see README.md of the platform
	)

	scaler := aio.AnalogSensorLinearScaler(0, 4095, 0, 1.8)

	adaptor := zero.NewAdaptor()
	ana0 := aio.NewAnalogSensorDriver(adaptor, inPin0, aio.WithSensorScaler(scaler))
	ana1 := aio.NewAnalogSensorDriver(adaptor, inPin1)

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			v0, err := ana0.Read()
			if err != nil {
				log.Println(err)
			}

			v1, err := ana1.Read()
			if err != nil {
				log.Println(err)
			}

			fmt.Printf("%s: %1.3f V, %s: %2.0f (%4.0f mV)\n", inPin0, v0, inPin1, v1, v1*inVoltageScale)
		})
	}

	robot := gobot.NewRobot("adcBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{ana0, ana1},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
