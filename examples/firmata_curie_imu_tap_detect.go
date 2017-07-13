// +build example
//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_curie_imu_tap_detect.go /dev/ttyACM0
*/

package main

import (
	"log"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/intel-iot/curie"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	led := gpio.NewLedDriver(firmataAdaptor, "13")
	imu := curie.NewIMUDriver(firmataAdaptor)

	work := func() {
		imu.On("Tap", func(data interface{}) {
			log.Println("Tap", data)
		})

		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})

		imu.EnableTapDetection(true)
	}

	robot := gobot.NewRobot("curieBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{imu, led},
		work,
	)

	robot.Start()
}
