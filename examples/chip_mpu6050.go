package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
  "github.com/hybridgroup/gobot/platforms/chip"
)

func main() {
	gbot := gobot.NewGobot()

	board := chip.NewChipAdaptor("chip")
	mpu6050 := i2c.NewMPU6050Driver(board, "mpu6050")

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
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

	gbot.AddRobot(robot)

	gbot.Start()
}
