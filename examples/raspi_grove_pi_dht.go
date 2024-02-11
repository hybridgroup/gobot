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

const (
	dhtPin        = "4"
	dhtModel      = 1 // white DHT22, 0 for the blue DHT11
	delayMillisec = 10
)

func main() {
	r := raspi.NewAdaptor()
	gp := i2c.NewGrovePiDriver(r)

	work := func() {
		gobot.Every(1*time.Second, func() {
			if temp, hum, err := gp.DHTRead(dhtPin, dhtModel, delayMillisec); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Temperature [°C]", temp)
				fmt.Println("rel. Humidity [%]", hum)
			}
		})
	}

	robot := gobot.NewRobot("dhtBot",
		[]gobot.Connection{r},
		[]gobot.Device{gp},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
