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

// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 2 (+5V), 6, 9, 14, 20 (GND)
// I2C1 Tinkerboard: 3 (SDA), 5 (SCL)
// PLATE: connected via pin header (pin 1..26)
func main() {
	board := tinkerboard.NewAdaptor()
	ada := i2c.NewAdafruit1109Driver(board, i2c.WithBus(1))

	work := func() {
		// set a custom character
		smiley := [8]byte{0, 0, 10, 0, 0, 17, 14, 0}
		if err := ada.CreateChar(0, smiley); err != nil {
			fmt.Println(err)
		}

		if err := ada.Clear(); err != nil {
			fmt.Println(err)
		}
		if err := ada.SetRGB(true, false, false); err != nil {
			fmt.Println(err)
		}
		if err := ada.Write("   Hello from   \n  Tinker Board "); err != nil {
			fmt.Println(err)
		}
		// add the custom character at the end of the string
		if err := ada.Write(string(byte(0))); err != nil {
			fmt.Println(err)
		}

		// after 1 sec. activate rotation
		direction := 1
		gobot.After(1*time.Second, func() {
			if err := ada.SetRGB(false, true, false); err != nil {
				fmt.Println(err)
			}
			gobot.Every(400*time.Millisecond, func() {
				if direction == 1 {
					if err := ada.ScrollLeft(); err != nil {
						fmt.Println(err)
					}
				}
				if direction == 2 {
					if err := ada.ScrollRight(); err != nil {
						fmt.Println(err)
					}
				}
			})
		})

		// after 7 sec. activate play with the buttons
		gobot.After(7*time.Second, func() {
			direction = 0
			time.Sleep(1 * time.Second)
			if err := ada.LeftToRight(); err != nil {
				fmt.Println(err)
			}
			if err := ada.Clear(); err != nil {
				fmt.Println(err)
			}
			if err := ada.SetRGB(false, false, true); err != nil {
				fmt.Println(err)
			}
			if err := ada.Write("Try the buttons!"); err != nil {
				fmt.Println(err)
			}

			gobot.Every(500*time.Millisecond, func() {
				if val, err := ada.SelectButton(); err != nil {
					fmt.Println(err)
				} else if val != 0 {
					if err := ada.Clear(); err != nil {
						fmt.Println(err)
					}
					if err := ada.Write("-Select Button-\nclear the screen"); err != nil {
						fmt.Println(err)
					}
					if err := ada.Blink(false); err != nil {
						fmt.Println(err)
					}
					direction = 0
				}
				if val, err := ada.UpButton(); err != nil {
					fmt.Println(err)
				} else if val != 0 {
					if err := ada.Clear(); err != nil {
						fmt.Println(err)
					}
					if err := ada.Write("  -Up Button-   \nset RGB to white"); err != nil {
						fmt.Println(err)
					}
					if err := ada.Blink(false); err != nil {
						fmt.Println(err)
					}
					if err := ada.SetRGB(true, true, true); err != nil {
						fmt.Println(err)
					}
					direction = 0
				}
				if val, err := ada.DownButton(); err != nil {
					fmt.Println(err)
				} else if val != 0 {
					if err := ada.Clear(); err != nil {
						fmt.Println(err)
					}
					if err := ada.Write(" -Down Button-  \nset blink on"); err != nil {
						fmt.Println(err)
					}
					if err := ada.Blink(true); err != nil {
						fmt.Println(err)
					}
					direction = 0
				}
				if val, err := ada.LeftButton(); err != nil {
					fmt.Println(err)
				} else if val != 0 {
					if err := ada.Clear(); err != nil {
						fmt.Println(err)
					}
					if err := ada.Write("   -Left Button-  \nrotate left"); err != nil {
						fmt.Println(err)
					}
					if err := ada.Blink(false); err != nil {
						fmt.Println(err)
					}
					direction = 1
				}
				if val, err := ada.RightButton(); err != nil {
					fmt.Println(err)
				} else if val != 0 {
					if err := ada.Clear(); err != nil {
						fmt.Println(err)
					}
					if err := ada.Write("   -Right Button-  \nrotate right"); err != nil {
						fmt.Println(err)
					}
					if err := ada.Blink(false); err != nil {
						fmt.Println(err)
					}
					direction = 2
				}
			})
		})
	}

	robot := gobot.NewRobot("adaBot",
		[]gobot.Connection{board},
		[]gobot.Device{ada},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
