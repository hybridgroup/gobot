package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/drivers/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

func main() {
	e := edison.NewAdaptor()
	touch := gpio.NewGroveTouchDriver(e, "2")

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

	robot.Start()
}
