package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/drivers/gpio"
	"github.com/hybridgroup/gobot/platforms/firmata"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	led1 := gpio.NewLedDriver(firmataAdaptor, "3")
	led2 := gpio.NewLedDriver(firmataAdaptor, "4")
	button := gpio.NewButtonDriver(firmataAdaptor, "2")
	sensor := gpio.NewAnalogSensorDriver(firmataAdaptor, "0")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led1.Toggle()
		})
		gobot.Every(2*time.Second, func() {
			led2.Toggle()
		})
		button.On(gpio.ButtonPush, func(data interface{}) {
			led2.On()
		})
		button.On(gpio.ButtonRelease, func(data interface{}) {
			led2.Off()
		})
		sensor.On(gpio.Data, func(data interface{}) {
			fmt.Println("sensor", data)
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led1, led2, button, sensor},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
