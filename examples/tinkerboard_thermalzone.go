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
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

// Wiring: no wiring needed
func main() {
	adaptor := tinkerboard.NewAdaptor()
	therm0 := aio.NewThermalZoneDriver(adaptor, "thermal_zone0")
	therm1 := aio.NewThermalZoneDriver(adaptor, "thermal_zone1", aio.WithFahrenheit())

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

			fmt.Printf("Zone 0: %2.3f °C, Zone 1: %2.3f °F\n", t0, t1)
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
