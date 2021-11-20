// +build example
//
// Do not build by default.

/*
 How to run

	go run examples/firmata_hmc8553l.go
*/

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	raspi := raspi.NewAdaptor()
	hmc8553l := i2c.NewHMC8553LDriver(raspi)

	work := func() {
		gobot.Every(200*time.Millisecond, func() {

			// get heading in radians, to convert to degrees multiply by 180/math.Pi
			heading, _ := hmc8553l.Heading()
			fmt.Println("Heading", heading)

			// read the raw data from the device, this is useful for calibration
			x, y, z, _ := hmc8553l.ReadRawData()
			fmt.Println(x, y, z)
		})
	}

	robot := gobot.NewRobot("hmc8553LBot",
		[]gobot.Connection{raspi},
		[]gobot.Device{hmc8553l},
		work,
	)

	robot.Start()
}
