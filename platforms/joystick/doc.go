/*
Package joystick provides the Gobot adaptor and drivers for game controllers that are compatible with SDL.

Installing:

This package requires `sdl2` to be installed on your system
Then install package with:

	go get github.com/hybridgroup/gobot/platforms/joystick

Example:

	package main

	import (
		"fmt"

		"github.com/hybridgroup/gobot"
		"github.com/hybridgroup/gobot/platforms/joystick"
	)

	func main() {
		gbot := gobot.NewGobot()

		joystickAdaptor := joystick.NewJoystickAdaptor("ps3")
		joystick := joystick.NewJoystickDriver(joystickAdaptor,
			"ps3",
			"./platforms/joystick/configs/dualshock3.json",
		)

		work := func() {
			gobot.On(joystick.Event("square_press"), func(data interface{}) {
				fmt.Println("square_press")
			})
			gobot.On(joystick.Event("square_release"), func(data interface{}) {
				fmt.Println("square_release")
			})
			gobot.On(joystick.Event("triangle_press"), func(data interface{}) {
				fmt.Println("triangle_press")
			})
			gobot.On(joystick.Event("triangle_release"), func(data interface{}) {
				fmt.Println("triangle_release")
			})
			gobot.On(joystick.Event("left_x"), func(data interface{}) {
				fmt.Println("left_x", data)
			})
			gobot.On(joystick.Event("left_y"), func(data interface{}) {
				fmt.Println("left_y", data)
			})
			gobot.On(joystick.Event("right_x"), func(data interface{}) {
				fmt.Println("right_x", data)
			})
			gobot.On(joystick.Event("right_y"), func(data interface{}) {
				fmt.Println("right_y", data)
			})
		}

		robot := gobot.NewRobot("joystickBot",
			[]gobot.Connection{joystickAdaptor},
			[]gobot.Device{joystick},
			work,
		)

		gbot.AddRobot(robot)

		gbot.Start()
	}

For further information refer to joystick README:
https://github.com/hybridgroup/gobot/blob/master/platforms/joystick/README.md
*/
package joystick
