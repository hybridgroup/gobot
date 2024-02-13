/*
Package joystick provides the Gobot adaptor and drivers for game controllers and joysticks.

Installing:

	Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Example:

	package main

	import (
		"fmt"

		"gobot.io/x/gobot/v2"
		"gobot.io/x/gobot/v2/platforms/joystick"
	)

	func main() {
		joystickAdaptor := joystick.NewAdaptor("0")
		joystick := joystick.NewDriver(joystickAdaptor, "dualshock3")

		work := func() {
			_ = joystick.On(joystick.Event("square_press"), func(data interface{}) {
				fmt.Println("square_press")
			})
			_ = joystick.On(joystick.Event("square_release"), func(data interface{}) {
				fmt.Println("square_release")
			})
			_ = joystick.On(joystick.Event("triangle_press"), func(data interface{}) {
				fmt.Println("triangle_press")
			})
			_ = joystick.On(joystick.Event("triangle_release"), func(data interface{}) {
				fmt.Println("triangle_release")
			})
			_ = joystick.On(joystick.Event("left_x"), func(data interface{}) {
				fmt.Println("left_x", data)
			})
			_ = joystick.On(joystick.Event("left_y"), func(data interface{}) {
				fmt.Println("left_y", data)
			})
			_ = joystick.On(joystick.Event("right_x"), func(data interface{}) {
				fmt.Println("right_x", data)
			})
			_ = joystick.On(joystick.Event("right_y"), func(data interface{}) {
				fmt.Println("right_y", data)
			})
		}

		robot := gobot.NewRobot("joystickBot",
			[]gobot.Connection{joystickAdaptor},
			[]gobot.Device{joystick},
			work,
		)

		if err := robot.Start(); err != nil {
			panic(err)
		}
	}

For further information refer to joystick README:
https://github.com/hybridgroup/gobot/blob/release/platforms/joystick/README.md
*/
package joystick // import "gobot.io/x/gobot/v2/platforms/joystick"
