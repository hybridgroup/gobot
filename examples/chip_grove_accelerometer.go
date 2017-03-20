// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/chip"
)

func main() {
	board := chip.NewAdaptor()
	accel := i2c.NewGroveAccelerometerDriver(board)

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			if x, y, z, err := accel.XYZ(); err == nil {
				fmt.Println(x, y, z)
				fmt.Println(accel.Acceleration(x, y, z))
			} else {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("accelBot",
		[]gobot.Connection{board},
		[]gobot.Device{accel},
		work,
	)

	robot.Start()
}
