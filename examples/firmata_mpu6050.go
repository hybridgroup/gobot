package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/drivers/i2c"
	"github.com/hybridgroup/gobot/platforms/firmata"
)

func main() {
	gbot := gobot.NewMaster()

	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	mpu6050 := i2c.NewMPU6050Driver(firmataAdaptor)

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			fmt.Println("Accelerometer", mpu6050.Accelerometer)
			fmt.Println("Gyroscope", mpu6050.Gyroscope)
			fmt.Println("Temperature", mpu6050.Temperature)
		})
	}

	robot := gobot.NewRobot("mpu6050Bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{mpu6050},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
