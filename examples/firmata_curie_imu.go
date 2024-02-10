//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_curie_imu.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
	"gobot.io/x/gobot/v2/platforms/intel-iot/curie"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	led := gpio.NewLedDriver(firmataAdaptor, "13")
	imu := curie.NewIMUDriver(firmataAdaptor)

	work := func() {
		_ = imu.On("Accelerometer", func(data interface{}) {
			log.Println("Accelerometer", data)
		})

		_ = imu.On("Gyroscope", func(data interface{}) {
			log.Println("Gyroscope", data)
		})

		_ = imu.On("Temperature", func(data interface{}) {
			log.Println("Temperature", data)
		})

		_ = imu.On("Motion", func(data interface{}) {
			log.Println("Motion", data)
		})

		gobot.Every(1*time.Second, func() {
			if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
		})

		gobot.Every(100*time.Millisecond, func() {
			if err := imu.ReadAccelerometer(); err != nil {
				fmt.Println(err)
			}
			if err := imu.ReadGyroscope(); err != nil {
				fmt.Println(err)
			}
			if err := imu.ReadTemperature(); err != nil {
				fmt.Println(err)
			}
			if err := imu.ReadMotion(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("curieBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{imu, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
