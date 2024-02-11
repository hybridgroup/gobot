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
	"gobot.io/x/gobot/v2/platforms/chip"
)

func main() {
	board := chip.NewAdaptor()
	luxSensor := i2c.NewTSL2561Driver(board, i2c.WithTSL2561Gain16X)

	work := func() {
		gobot.Every(1*time.Second, func() {
			broadband, ir, err := luxSensor.GetLuminocity()

			if err != nil {
				fmt.Println("Err:", err)
			} else {
				light := luxSensor.CalculateLux(broadband, ir)
				fmt.Printf("BB: %v, IR: %v, Lux: %v\n", broadband, ir, light)
			}
		})
	}

	robot := gobot.NewRobot("tsl2561Bot",
		[]gobot.Connection{board},
		[]gobot.Device{luxSensor},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
