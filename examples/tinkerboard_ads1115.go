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
// ADS1115 plate: VDD (2.0..5.5V), GND, SDL, SDA, A0..A3 (input voltage 0..5V)
func main() {
	const voltage = 5.0 // we will be able to read values of at least 5V (for all channels)

	a := tinkerboard.NewAdaptor()
	ads1115 := i2c.NewADS1115Driver(a, i2c.WithADS1x15BestGainForVoltage(voltage))

	work := func() {
		var a0, a1, a2, a3 float64
		var err error
		gobot.Every(2*time.Second, func() {
			a0, err = ads1115.ReadWithDefaults(0)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("A0", a0)
			}
			a1, err = ads1115.ReadWithDefaults(1)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("A1", a1)
			}
			a2, err = ads1115.ReadWithDefaults(2)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("A2", a2)
			}
			a3, err = ads1115.ReadWithDefaults(3)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("A3", a3)
			}
			if v, err := ads1115.ReadDifferenceWithDefaults(0); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("A0-A1 %f, (e: %f)\n", v, v-(a0-a1))
			}
			if v, err := ads1115.ReadDifferenceWithDefaults(1); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("A0-A3 %f, (e: %f)\n", v, v-(a0-a3))
			}
			if v, err := ads1115.ReadDifferenceWithDefaults(2); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("A1-A3 %f, (e: %f)\n", v, v-(a1-a3))
			}
			if v, err := ads1115.ReadDifferenceWithDefaults(3); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("A2-A3 %f, (e: %f)\n", v, v-(a2-a3))
			}
			fmt.Println("-------------")
		})
	}

	robot := gobot.NewRobot("ads1115bot",
		[]gobot.Connection{a},
		[]gobot.Device{ads1115},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
