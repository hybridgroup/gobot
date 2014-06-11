package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("myFirmata", "/dev/ttyACM0")

	button := gpio.NewButtonDriver(firmataAdaptor, "myButton", "2")
	led := gpio.NewLedDriver(firmataAdaptor, "myLed", "13")

	work := func() {
		gobot.On(button.Events["push"], func(data interface{}) {
			led.On()
		})
		gobot.On(button.Events["release"], func(data interface{}) {
			led.Off()
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("buttonBot", []gobot.Connection{firmataAdaptor}, []gobot.Device{button, led}, work),
	)

	gbot.Start()
}
