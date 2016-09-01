package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	board := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	led := gpio.NewRgbLedDriver(board, "led", "3", "5", "6")

	work := func() {
		gobot.Every(1*time.Second, func() {
			r := uint8(gobot.Rand(255))
			g := uint8(gobot.Rand(255))
			b := uint8(gobot.Rand(255))
			led.SetRGB(r, g, b)
		})
	}

	robot := gobot.NewRobot("rgbBot",
		[]gobot.Connection{board},
		[]gobot.Device{led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
