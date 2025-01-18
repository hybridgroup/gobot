//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

func main() {
	board := tinkerboard.NewAdaptor()
	screen := i2c.NewGroveLcdDriver(board)

	work := func() {
		if err := screen.Write("Hello from"); err != nil {
			fmt.Println(err)
		}

		if err := screen.SetRGB(255, 0, 0); err != nil {
			fmt.Println(err)
		}

		gobot.After(5*time.Second, func() {
			if err := screen.Clear(); err != nil {
				fmt.Println(err)
			}
			if err := screen.Home(); err != nil {
				fmt.Println(err)
			}
			if err := screen.SetRGB(0, 255, 0); err != nil {
				fmt.Println(err)
			}
			// set a custom character in the first position
			if err := screen.SetCustomChar(0, i2c.CustomLCDChars["smiley"]); err != nil {
				fmt.Println(err)
			}
			// add the custom character at the end of the string
			if err := screen.Write("\nTinker Board " + string(byte(0))); err != nil {
				fmt.Println(err)
			}
			gobot.Every(500*time.Millisecond, func() {
				if err := screen.Scroll(false); err != nil {
					fmt.Println(err)
				}
			})
		})

		if err := screen.Home(); err != nil {
			fmt.Println(err)
		}
		time.Sleep(1 * time.Second)
		if err := screen.SetRGB(0, 0, 255); err != nil {
			fmt.Println(err)
		}
	}

	robot := gobot.NewRobot("screenBot",
		[]gobot.Connection{board},
		[]gobot.Device{screen},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
