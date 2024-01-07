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
	"gobot.io/x/gobot/v2/platforms/raspi"
)

// Wiring: no wiring needed
func main() {
	adaptor := raspi.NewAdaptor()
	therm0C := aio.NewThermalZoneDriver(adaptor, "thermal_zone0")
	therm0F := aio.NewThermalZoneDriver(adaptor, "thermal_zone0", aio.WithFahrenheit())

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			t0C, err := therm0C.Read()
			if err != nil {
				log.Println(err)
			}

			t0F, err := therm0F.Read()
			if err != nil {
				log.Println(err)
			}

			fmt.Printf("Zone 0: %2.3f °C, %2.3f °F\n", t0C, t0F)
		})
	}

	robot := gobot.NewRobot("thermalBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{therm0C, therm0F},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
