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
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	sht3x := i2c.NewSHT3xDriver(r)

	work := func() {
		sht3x.Units = "F"
		if err := sht3x.Start(); err != nil {
			fmt.Println(err)
		}
		sn, err := sht3x.SerialNumber()
		fmt.Printf("Serial Number: 0x%08x, err: %v\n", sn, err)

		gobot.Every(5*time.Second, func() {
			temp, rh, err := sht3x.Sample()
			fmt.Printf("Temp: %f F, Relative Humidity: %f, err: %v\n", temp, rh, err)
		})
	}

	robot := gobot.NewRobot("SHT3xbot",
		[]gobot.Connection{r},
		[]gobot.Device{sht3x},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
