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
	a := raspi.NewAdaptor()
	// Use the gain to be able to read values of at least 5V (for all channels)
	ads1015 := i2c.NewADS1015Driver(a, i2c.WithADS1x15BestGainForVoltage(5.0))

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			v, _ := ads1015.ReadWithDefaults(0)
			fmt.Println("A0", v)
		})
	}

	robot := gobot.NewRobot("ads1015bot",
		[]gobot.Connection{a},
		[]gobot.Device{ads1015},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
