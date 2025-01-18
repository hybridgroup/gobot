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
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 6, 9, 14, 20 (GND)
// I2C1 Tinkerboard: 3 (SDA), 5 (SCL)
// HMC5883L plate: VCC, GND, SDL, SDA
func main() {
	a := tinkerboard.NewAdaptor()
	hmc := i2c.NewHMC5883LDriver(a)

	work := func() {
		var x, y, z, h float64
		var err error

		gobot.Every(1000*time.Millisecond, func() {
			if x, y, z, err = hmc.Read(); err != nil {
				fmt.Println(err)
			} else {
				log.Printf("Xn: %f, Yn: %f, Zn: %f", x, y, z)
			}
			if h, err = hmc.Heading(); err != nil {
				fmt.Println(err)
			} else {
				log.Printf("Heading: %.1f", h)
			}
		})
	}

	robot := gobot.NewRobot("hcBot",
		[]gobot.Connection{a},
		[]gobot.Device{hmc},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
