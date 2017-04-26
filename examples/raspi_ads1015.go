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
	a := raspi.NewAdaptor()
	ads1015 := i2c.NewADS1015Driver(a)
	// Adjust the gain to be able to read values of at least 5V
	ads1015.DefaultGain, _ = ads1015.BestGainForVoltage(5.0)

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

	robot.Start()
}
