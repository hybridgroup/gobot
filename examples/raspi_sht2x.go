// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	sht2x := i2c.NewSHT2xDriver(r)

	work := func() {
		gobot.Every(1*time.Second, func() {
			t, _ := sht2x.Temperature()
			fmt.Printf("Temperature: %v\n", t)

			h, _ := sht2x.Humidity()
			fmt.Printf("Humidity: %v\n", h)
		})
	}

	robot := gobot.NewRobot("SHT2xbot",
		[]gobot.Connection{r},
		[]gobot.Device{sht2x},
		work,
	)

	robot.Start()
}
