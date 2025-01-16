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
	"gobot.io/x/gobot/v2/platforms/friendlyelec/nanopct6"
)

// Wiring: no wiring needed
func main() {
	adaptor := nanopct6.NewAdaptor()
	therm0 := aio.NewThermalZoneDriver(adaptor, "soc_thermal")
	therm1 := aio.NewThermalZoneDriver(adaptor, "npu_thermal", aio.WithFahrenheit())

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			t0, err := therm0.Read()
			if err != nil {
				log.Println(err)
			}

			t1, err := therm1.Read()
			if err != nil {
				log.Println(err)
			}

			fmt.Printf("SOC: %2.3f °C, NPU: %2.3f °F\n", t0, t1)
		})
	}

	robot := gobot.NewRobot("thermalBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{therm0, therm1},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
