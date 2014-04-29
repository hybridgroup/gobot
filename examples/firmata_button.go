package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/core/api"
	"github.com/hybridgroup/gobot/core/robot"
	"github.com/hybridgroup/gobot/core/utils"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func main() {
	gbot := gobot.NewGobot()
	gbot.Api = api.NewApi()

	firmataAdaptor := firmata.NewFirmataAdaptor("myFirmata", "/dev/ttyACM0")

	button := gpio.NewButtonDriver(firmataAdaptor, "myButton", "2")
	led := gpio.NewLedDriver(firmataAdaptor, "myLed", "13")

	work := func() {
		utils.Every((1 * time.Second), func() {
			led.Toggle()
		})
		utils.On(button.Events["push"], func(data interface{}) {
			led.On()
		})
		utils.On(button.Events["release"], func(data interface{}) {
			led.Off()
		})
	}

	gbot.Robots = append(gbot.Robots,
		robot.NewRobot("name", []robot.Connection{firmataAdaptor}, []robot.Device{button, led}, work),
	)

	gbot.Start()
}
