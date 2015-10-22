package main

import (
	"time"
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	led1 := gpio.NewLedDriver(firmataAdaptor, "led1", "3")
	led2 := gpio.NewLedDriver(firmataAdaptor, "led2", "4")
	button := gpio.NewButtonDriver(firmataAdaptor, "button", "2")
	sensor := gpio.NewAnalogSensorDriver(firmataAdaptor, "sensor", "0")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led1.Toggle()
		})
		gobot.Every(2*time.Second, func() {
			led2.Toggle()
		})
		gobot.On(button.Event("push"), func(data interface{}) {
			led2.On()
		})
		gobot.On(button.Event("release"), func(data interface{}) {
			led2.Off()
		})
		gobot.On(sensor.Event("data"), func(data interface{}) {
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
