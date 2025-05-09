//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 2(+5V), 6, 9, 14, 20 (GND)
// I2C1 Tinkerboard: 3 (SDA-ws), 5 (SCL-gn)
// ADXL345 plate: VCC, GND, SDL, SDA
func main() {
	a := tinkerboard.NewAdaptor()
	adxl := i2c.NewADXL345Driver(a)

	work := func() {
		gobot.Every(1000*time.Millisecond, func() {
			if x, y, z, err := adxl.XYZ(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("x: %.7f | y: %.7f | z: %.7f \n", x, y, z)
			}
		})
	}

	robot := gobot.NewRobot("mpBot",
		[]gobot.Connection{a},
		[]gobot.Device{adxl},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
