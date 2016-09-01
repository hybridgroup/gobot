package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"

	"github.com/hybridgroup/gobot/api"
)

func main() {
	gbot := gobot.NewGobot()
	api.NewAPI(gbot).Start()

	e := edison.NewEdisonAdaptor("edison")

	button := gpio.NewButtonDriver(e, "myButton", "2")
	led := gpio.NewLedDriver(e, "myLed", "4")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			led.On()
		})
		button.On(gpio.ButtonRelease, func(data interface{}) {
			led.Off()
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{e},
		[]gobot.Device{led, button},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
