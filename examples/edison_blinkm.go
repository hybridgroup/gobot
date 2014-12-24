package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

func main() {
	gbot := gobot.NewGobot()

	e := edison.NewEdisonAdaptor("edison")
	blinkm := i2c.NewBlinkMDriver(e, "blinkm")

	work := func() {
		gobot.Every(3*time.Second, func() {
			r := byte(gobot.Rand(255))
			g := byte(gobot.Rand(255))
			b := byte(gobot.Rand(255))
			blinkm.Rgb(r, g, b)
			color, _ := blinkm.Color()
			fmt.Println("color", color)
		})
	}

	robot := gobot.NewRobot("blinkmBot",
		[]gobot.Connection{e},
		[]gobot.Device{blinkm},
		work,
	)

	gbot.AddRobot(robot)
	gbot.Start()
}
