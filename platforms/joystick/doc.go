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
			joystick.On(joystick.Event("square_press"), func(data interface{}) {
				fmt.Println("square_press")
			})
			joystick.On(joystick.Event("square_release"), func(data interface{}) {
				fmt.Println("square_release")
			})
			joystick.On(joystick.Event("triangle_press"), func(data interface{}) {
				fmt.Println("triangle_press")
			})
			joystick.On(joystick.Event("triangle_release"), func(data interface{}) {
				fmt.Println("triangle_release")
			})
			joystick.On(joystick.Event("left_x"), func(data interface{}) {
				fmt.Println("left_x", data)
			})
			joystick.On(joystick.Event("left_y"), func(data interface{}) {
				fmt.Println("left_y", data)
			})
			joystick.On(joystick.Event("right_x"), func(data interface{}) {
				fmt.Println("right_x", data)
			})
			joystick.On(joystick.Event("right_y"), func(data interface{}) {
				fmt.Println("right_y", data)
			})
		}

		robot := gobot.NewRobot("joystickBot",
			[]gobot.Connection{joystickAdaptor},
			[]gobot.Device{joystick},
			work,
		)

		robot.Start()
	}

For further information refer to joystick README:
https://github.com/hybridgroup/gobot/blob/master/platforms/joystick/README.md
*/
package joystick // import "gobot.io/x/gobot/v2/platforms/joystick"
