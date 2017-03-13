// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")

	sensor := gpio.NewPIRMotionDriver(firmataAdaptor, "5")
	led := gpio.NewLedDriver(firmataAdaptor, "13")

	work := func() {
		sensor.On(gpio.MotionDetected, func(data interface{}) {
			fmt.Println(gpio.MotionDetected)
			led.On()
		})
		sensor.On(gpio.MotionStopped, func(data interface{}) {
			fmt.Println(gpio.MotionStopped)
			led.Off()
		})
	}

	robot := gobot.NewRobot("motionBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{sensor, led},
		work,
	)

	robot.Start()
}
