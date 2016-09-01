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
		touch.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("On!")
		})

		touch.On(gpio.ButtonRelease, func(data interface{}) {
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
