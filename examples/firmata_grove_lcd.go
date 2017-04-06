// +build example
//
// Do not build by default.

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	board := firmata.NewAdaptor(os.Args[1])
	screen := i2c.NewGroveLcdDriver(board)

	work := func() {
		screen.Write("hello")

		screen.SetRGB(255, 0, 0)

		gobot.After(5*time.Second, func() {
			screen.Clear()
			screen.Home()
			screen.SetRGB(0, 255, 0)
			// set a custom character in the first position
			screen.SetCustomChar(0, i2c.CustomLCDChars["smiley"])
			// add the custom character at the end of the string
			screen.Write("goodbye\nhave a nice day " + string(byte(0)))
			gobot.Every(500*time.Millisecond, func() {
				screen.Scroll(false)
			})
		})

		screen.Home()
		time.Sleep(1 * time.Second)
		screen.SetRGB(0, 0, 255)
	}

	robot := gobot.NewRobot("screenBot",
		[]gobot.Connection{board},
		[]gobot.Device{screen},
		work,
	)

	robot.Start()
}
