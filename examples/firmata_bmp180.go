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
	bmp180 := i2c.NewBMP180Driver(firmataAdaptor)

	work := func() {
		gobot.Every(1*time.Second, func() {
			t, _ := bmp180.Temperature()
			fmt.Println("Temperature", t)

			p, _ := bmp180.Pressure()
			fmt.Println("Pressure", p)
		})
	}

	robot := gobot.NewRobot("bmp180bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{bmp180},
		work,
	)

	robot.Start()
}
