// +build example
//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_adxl345.go /dev/ttyACM0
*/

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
	adxl345 := i2c.NewADXL345Driver(firmataAdaptor)

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			x, y, z, _ := adxl345.XYZ()

			fmt.Printf("x: %.7f | y: %.7f | z: %.7f \n", x, y, z)
		})
	}

	robot := gobot.NewRobot("adxl345Bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{adxl345},
		work,
	)

	robot.Start()
}
