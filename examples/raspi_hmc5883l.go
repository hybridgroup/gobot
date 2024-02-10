//go:build example
// +build example

//
// Do not build by default.

/*
 How to run

	go run examples/raspi_hmc5883l.go
*/

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	raspi := raspi.NewAdaptor()
	hmc5883l := i2c.NewHMC5883LDriver(raspi)

	work := func() {
		gobot.Every(200*time.Millisecond, func() {
			// get heading in radians, to convert to degrees multiply by 180/math.Pi
			heading, _ := hmc5883l.Heading()
			fmt.Println("Heading", heading)

			// read the data in Gauss
			x, y, z, _ := hmc5883l.Read()
			fmt.Println(x, y, z)
		})
	}

	robot := gobot.NewRobot("hmc5883LBot",
		[]gobot.Connection{raspi},
		[]gobot.Device{hmc5883l},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
