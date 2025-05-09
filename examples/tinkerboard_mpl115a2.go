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
// MPL115A2 plate: VDD (2.375..5.5V), GND, SDL, SDA
func main() {
	board := tinkerboard.NewAdaptor()
	mpl115a2 := i2c.NewMPL115A2Driver(board)

	work := func() {
		gobot.Every(2*time.Second, func() {
			if press, err := mpl115a2.Pressure(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Pressure [kPa]", press)
			}

			if temp, err := mpl115a2.Temperature(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Temperature [°C]", temp)
			}

			fmt.Println("-------------")
		})
	}

	robot := gobot.NewRobot("mpl115Bot",
		[]gobot.Connection{board},
		[]gobot.Device{mpl115a2},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
