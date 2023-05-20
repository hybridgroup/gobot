//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass the BLE address or BLE name as first param:

	go run examples/ble_firmata_curie_imu.go FIRMATA

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"log"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
	"gobot.io/x/gobot/v2/platforms/intel-iot/curie"
)

func main() {
	var bleAdaptor interface{} = firmata.NewBLEAdaptor(os.Args[1])
	firmataAdaptor := bleAdaptor.(*firmata.Adaptor)
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

		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})

		gobot.Every(100*time.Millisecond, func() {
			imu.ReadAccelerometer()
			imu.ReadGyroscope()
			imu.ReadTemperature()
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led, imu},
		work,
	)

	robot.Start()
}
