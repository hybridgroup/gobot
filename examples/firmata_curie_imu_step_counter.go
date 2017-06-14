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
		imu.On("Steps", func(data interface{}) {
			log.Println("Steps", data)
		})

		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})

		imu.EnableStepCounter(true)
	}

	robot := gobot.NewRobot("curieBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{imu, led},
		work,
	)

	robot.Start()
}
