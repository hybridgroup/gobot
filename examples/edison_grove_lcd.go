package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

func main() {
	gbot := gobot.NewGobot()

	board := edison.NewEdisonAdaptor("edison")
	screen := i2c.NewGroveLcdDriver(board, "screen")

	work := func() {
		screen.Write("hello")

		screen.SetRGB(255, 0, 0)

		gobot.After(5*time.Second, func() {
			screen.Clear()
			screen.Home()
			screen.SetRGB(0, 255, 0)
			screen.Write("goodbye")
		})

		screen.Home()
		<-time.After(1 * time.Second)
		screen.SetRGB(0, 0, 255)
	}

	robot := gobot.NewRobot("screenBot",
		[]gobot.Connection{board},
		[]gobot.Device{screen},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
