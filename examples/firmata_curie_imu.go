// +build example
//
// Do not build by default.

package main

import (
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/intel-iot/curie"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "13")
	imu := curie.NewIMUDriver(firmataAdaptor)

	work := func() {
		imu.On("Accelerometer", func(data interface{}) {
			log.Println("Accelerometer", data)
		})

		imu.On("Gyroscope", func(data interface{}) {
			log.Println("Gyroscope", data)
		})

		imu.On("Temperature", func(data interface{}) {
			log.Println("Temperature", data)
		})

		imu.On("Motion", func(data interface{}) {
			log.Println("Motion", data)
		})

		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})

		gobot.Every(100*time.Millisecond, func() {
			imu.ReadAccelerometer()
			imu.ReadGyroscope()
			imu.ReadTemperature()
			imu.ReadMotion()
		})
	}

	robot := gobot.NewRobot("blinkmBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{imu, led},
		work,
	)

	robot.Start()
}
