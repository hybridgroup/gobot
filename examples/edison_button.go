package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

func main() {
	gbot := gobot.NewGobot()

	e := edison.NewEdisonAdaptor("edison")

	button := gpio.NewButtonDriver(e, "myButton", "2")
	led := gpio.NewLedDriver(e, "myLed", "7")

	work := func() {
		gobot.On(button.Event("push"), func(data interface{}) {
			led.On()
		})
		gobot.On(button.Event("release"), func(data interface{}) {
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
