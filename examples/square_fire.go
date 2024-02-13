//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/intel-iot/edison"
)

func main() {
	manager := gobot.NewManager()
	a := api.NewAPI(manager)
	a.Start()

	board := edison.NewAdaptor()
	red := gpio.NewLedDriver(board, "3")
	green := gpio.NewLedDriver(board, "5")
	blue := gpio.NewLedDriver(board, "6")

	button := gpio.NewButtonDriver(board, "7")

	enabled := true
	work := func() {
		if err := red.Brightness(0xff); err != nil {
			fmt.Println(err)
		}
		if err := green.Brightness(0x00); err != nil {
			fmt.Println(err)
		}
		if err := blue.Brightness(0x00); err != nil {
			fmt.Println(err)
		}

		flash := false
		on := true

		gobot.Every(50*time.Millisecond, func() {
			if enabled {
				if flash {
					if on {
						if err := red.Brightness(0x00); err != nil {
							fmt.Println(err)
						}
						if err := green.Brightness(0xff); err != nil {
							fmt.Println(err)
						}
						if err := blue.Brightness(0x00); err != nil {
							fmt.Println(err)
						}
						on = false
					} else {
						if err := red.Brightness(0x00); err != nil {
							fmt.Println(err)
						}
						if err := green.Brightness(0x00); err != nil {
							fmt.Println(err)
						}
						if err := blue.Brightness(0xff); err != nil {
							fmt.Println(err)
						}
						on = true
					}
				}
			}
		})

		_ = button.On(gpio.ButtonPush, func(data interface{}) {
			flash = true
		})

		_ = button.On(gpio.ButtonRelease, func(data interface{}) {
			flash = false
			if err := red.Brightness(0x00); err != nil {
				fmt.Println(err)
			}
			if err := green.Brightness(0x00); err != nil {
				fmt.Println(err)
			}
			if err := blue.Brightness(0xff); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot(
		"square of fire",
		[]gobot.Connection{board},
		[]gobot.Device{red, green, blue, button},
		work,
	)

	robot.AddCommand("enable", func(params map[string]interface{}) interface{} {
		enabled = !enabled
		return enabled
	})

	manager.AddRobot(robot)

	if err := manager.Start(); err != nil {
		panic(err)
	}
}
