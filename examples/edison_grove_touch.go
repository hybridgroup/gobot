package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

func main() {
	gbot := gobot.NewGobot()

	e := edison.NewEdisonAdaptor("edison")
	touch := gpio.NewGroveTouchDriver(e, "touch", "2")

	work := func() {
		gobot.On(touch.Event(gpio.Push), func(data interface{}) {
			fmt.Println("On!")
		})

		gobot.On(touch.Event(gpio.Release), func(data interface{}) {
			fmt.Println("Off!")
		})

	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{e},
		[]gobot.Device{touch},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
