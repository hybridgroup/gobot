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
	"gobot.io/x/gobot/v2/drivers/onewire"
	"gobot.io/x/gobot/v2/platforms/friendlyelec/nanopct6"
)

// Preparation: see /gobot/system/ONEWIRE.md
//
// Wiring:
// PWR    : 1, 17 (+3.3V, VCC), 6, 9, 14, 20, 25, 30, 34, 39 (GND)
// 1-wire : 16 (DQ) - resistor to VCC, ~1.5kOhm ... 5kOhm
// DS18B20: 1 (GND), 2 (DQ), 3 (VDD, +3 ... 5.5V) for local power mode
func main() {
	adaptor := nanopct6.NewAdaptor()
	// resolution change not supported by all devices
	temp0 := onewire.NewDS18B20Driver(adaptor, 0xde5e710a6461, onewire.WithResolution(12))
	temp1 := onewire.NewDS18B20Driver(adaptor, 0x1e40710a6461, onewire.WithFahrenheit(), onewire.WithConversionTime(500))

	work := func() {
		time0, err := temp0.ConversionTime()
		if err != nil {
			log.Printf("Err CT0: %v\n", err)
		}
		res0, err := temp0.Resolution()
		if err != nil {
			log.Printf("Err R0: %v\n", err)
		}
		log.Printf("Conversion time @%d bit for Temp 0: %d ms\n", res0, time0)

		time1, err := temp1.ConversionTime()
		if err != nil {
			log.Printf("Err CT1: %v\n", err)
		}
		res1, err := temp1.Resolution()
		if err != nil {
			log.Printf("Err R1: %v\n", err)
		}
		log.Printf("Conversion time @%d bit for Temp 0: %d ms\n", res1, time1)

		gobot.Every(10*(time.Duration(time0))*time.Millisecond, func() {
			t0, err := temp0.Temperature()
			if err != nil {
				log.Printf("Err Temp 0: %v\n", err)
			}

			fmt.Printf("Temp 0: %2.1f °C\n", t0)
		})

		gobot.Every(10*(time.Duration(time1))*time.Millisecond, func() {
			t1, err := temp1.Temperature()
			if err != nil {
				log.Printf("Err Temp 1:  %v\n", err)
			}

			fmt.Printf("Temp 1: %2.3f °F\n", t1)
		})
	}

	robot := gobot.NewRobot("onewireBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{temp0, temp1},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
