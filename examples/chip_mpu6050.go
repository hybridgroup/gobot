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
	mpu6050 := i2c.NewMPU6050Driver(board)

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			mpu6050.GetData()

			fmt.Println("Accelerometer", mpu6050.Accelerometer)
			fmt.Println("Gyroscope", mpu6050.Gyroscope)
			fmt.Println("Temperature", mpu6050.Temperature)
		})
	}

	robot := gobot.NewRobot("mpu6050Bot",
		[]gobot.Connection{board},
		[]gobot.Device{mpu6050},
		work,
	)

	robot.Start()
}
