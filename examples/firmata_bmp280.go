//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	bmp280 := i2c.NewBMP280Driver(firmataAdaptor)

	work := func() {
		gobot.Every(1*time.Second, func() {
			t, _ := bmp280.Temperature()
			fmt.Println("Temperature", t)

			p, _ := bmp280.Pressure()
			fmt.Println("Pressure", p)

			a, _ := bmp280.Altitude()
			fmt.Println("Altitude", a)
		})
	}

	robot := gobot.NewRobot("bmp280bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{bmp280},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
