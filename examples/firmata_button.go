package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("myFirmata", "/dev/ttyACM0")

	button := gpio.NewButtonDriver(firmataAdaptor, "myButton", "2")
	led := gpio.NewLedDriver(firmataAdaptor, "myLed", "13")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
		gobot.On(button.Events["push"], func(data interface{}) {
			led.On()
		})
		gobot.On(button.Events["release"], func(data interface{}) {
			led.Off()
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("buttonBot", []robot.Connection{firmataAdaptor}, []robot.Device{button, led}, work),
	)

	gbot.Start()
}
