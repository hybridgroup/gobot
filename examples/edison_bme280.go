// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
)

func main() {
	a := edison.NewAdaptor()
	bme280 := i2c.NewBME280Driver(a, i2c.WithAddress(0x76))

	work := func() {
		gobot.Every(1*time.Second, func() {
			t, e := bme280.Temperature()
			fmt.Println("Temperature", t)
			if e != nil {
				fmt.Println(e)
			}

			p, e := bme280.Pressure()
			fmt.Println("Pressure", p)
			if e != nil {
				fmt.Println(e)
			}

			a, e := bme280.Altitude()
			fmt.Println("Altitude", a)
			if e != nil {
				fmt.Println(e)
			}

			h, e := bme280.Humidity()
			fmt.Println("Humidity", h)
			if e != nil {
				fmt.Println(e)
			}
		})
	}

	robot := gobot.NewRobot("bme280bot",
		[]gobot.Connection{a},
		[]gobot.Device{bme280},
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
