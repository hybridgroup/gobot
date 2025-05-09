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
// BMP280 plate: VIN (2.5..5V), GND, SDL, SDA, SDO to GND
func main() {
	a := tinkerboard.NewAdaptor()
	bmp280 := i2c.NewBMP280Driver(a, i2c.WithAddress(0x76),
		i2c.WithBME280PressureOversampling(0x05),
		i2c.WithBME280TemperatureOversampling(0x02),
		i2c.WithBME280IIRFilter(0x05))

	work := func() {
		gobot.Every(2*time.Second, func() {
			t, e := bmp280.Temperature()
			fmt.Println("Temperature [°C]", t)
			if e != nil {
				fmt.Println(e)
			}

			p, e := bmp280.Pressure()
			fmt.Println("Pressure [Pa]", p) // 100hPa = 1Pa
			if e != nil {
				fmt.Println(e)
			}

			a, e := bmp280.Altitude()
			fmt.Println("Altitude [m]", a)
			if e != nil {
				fmt.Println(e)
			}
			fmt.Println("-------------")
		})
	}

	robot := gobot.NewRobot("bmp280bot",
		[]gobot.Connection{a},
		[]gobot.Device{bmp280},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
