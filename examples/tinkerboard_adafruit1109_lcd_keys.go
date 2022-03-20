// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/tinkerboard"
)

func main() {
	// * I2C1: 3 (SDA), 5 (SCL) --> connected to plate
	board := tinkerboard.NewAdaptor()
	ada := i2c.NewAdafruit1109Driver(board, i2c.WithBus(1))

	work := func() {
		// set a custom character
		smiley := [8]byte{0, 0, 10, 0, 0, 17, 14, 0}
		ada.CreateChar(0, smiley)

		ada.Clear()
		ada.SetRGB(true, false, false)
		ada.Write("   Hello from   \n  Tinker Board ")
		// add the custom character at the end of the string
		ada.Write(string(byte(0)))

		// after 1 sec. activate rotation
		direction := 1
		gobot.After(1*time.Second, func() {
			ada.SetRGB(false, true, false)
			gobot.Every(400*time.Millisecond, func() {
				if direction == 1 {
					ada.ScrollLeft()
				}
				if direction == 2 {
					ada.ScrollRight()
				}
			})
		})

		// after 7 sec. activate play with the buttons
		gobot.After(7*time.Second, func() {
			direction = 0
			time.Sleep(1 * time.Second)
			ada.LeftToRight()
			ada.Clear()
			ada.SetRGB(false, false, true)
			ada.Write("Try the buttons!")

			gobot.Every(500*time.Millisecond, func() {
				if val, err := ada.SelectButton(); err != nil {
					fmt.Println(err)
				} else {
					if val != 0 {
						ada.Clear()
						ada.Write("-Select Button-\nclear the screen")
						ada.Blink(false)
						direction = 0
					}
				}
				if val, err := ada.UpButton(); err != nil {
					fmt.Println(err)
				} else {
					if val != 0 {
						ada.Clear()
						ada.Write("  -Up Button-   \nset RGB to white")
						ada.Blink(false)
						ada.SetRGB(true, true, true)
						direction = 0
					}
				}
				if val, err := ada.DownButton(); err != nil {
					fmt.Println(err)
				} else {
					if val != 0 {
						ada.Clear()
						ada.Write(" -Down Button-  \nset blink on")
						ada.Blink(true)
						direction = 0
					}
				}
				if val, err := ada.LeftButton(); err != nil {
					fmt.Println(err)
				} else {
					if val != 0 {
						ada.Clear()
						ada.Write("   -Left Button-  \nrotate left")
						ada.Blink(false)
						direction = 1
					}
				}
				if val, err := ada.RightButton(); err != nil {
					fmt.Println(err)
				} else {
					if val != 0 {
						ada.Clear()
						ada.Write("   -Right Button-  \nrotate right")
						ada.Blink(false)
						direction = 2
					}
				}
			})

		})

	}

	robot := gobot.NewRobot("adaBot",
		[]gobot.Connection{board},
		[]gobot.Device{ada},
		work,
	)

	robot.Start()
}
