// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	bme280 := i2c.NewBME280Driver(firmataAdaptor, i2c.WithAddress(0x76))

	work := func() {
		gobot.Every(1*time.Second, func() {
			//fmt.Println("Pressure", mpl115a2.Pressure())
			t, _ := bme280.Temperature()
			fmt.Println("Temperature", t)

			p, _ := bme280.Pressure()
			fmt.Println("Pressure", p)

			h, _ := bme280.Humidity()
			fmt.Println("Humidity", h)
		})
	}

	robot := gobot.NewRobot("bme280bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{bme280},
		work,
	)

	robot.Start()
}
