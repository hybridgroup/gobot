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
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 2(+5V), 6, 9, 14, 20 (GND)
// I2C1 Tinkerboard: 3 (SDA-ws), 5 (SCL-gn)
// MPU6050 plate: VCC, GND, SDL, SDA
func main() {
	a := tinkerboard.NewAdaptor()
	mpu6050 := i2c.NewMPU6050Driver(a)

	work := func() {
		var err error

		gobot.Every(1000*time.Millisecond, func() {
			if err = mpu6050.GetData(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Acc: %v Gyr: %v Temp: %v\n", mpu6050.Accelerometer, mpu6050.Gyroscope, mpu6050.Temperature)
			}
		})
	}

	robot := gobot.NewRobot("mpBot",
		[]gobot.Connection{a},
		[]gobot.Device{mpu6050},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
