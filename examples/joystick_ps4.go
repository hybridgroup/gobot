// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/joystick"
)

func main() {
	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor,
		"./platforms/joystick/configs/dualshock4.json",
	)

	work := func() {
		stick.On(joystick.SquarePress, func(data interface{}) {
			fmt.Println("square_press")
		})
		stick.On(joystick.SquareRelease, func(data interface{}) {
			fmt.Println("square_release")
		})
		stick.On(joystick.TrianglePress, func(data interface{}) {
			fmt.Println("triangle_press")
		})
		stick.On(joystick.TriangleRelease, func(data interface{}) {
			fmt.Println("triangle_release")
		})
		stick.On(joystick.CirclePress, func(data interface{}) {
			fmt.Println("circle_press")
		})
		stick.On(joystick.CircleRelease, func(data interface{}) {
			fmt.Println("circle_release")
		})
		stick.On(joystick.XPress, func(data interface{}) {
			fmt.Println("x_press")
		})
		stick.On(joystick.XRelease, func(data interface{}) {
			fmt.Println("x_release")
		})
		stick.On(joystick.LeftX, func(data interface{}) {
			fmt.Println("left_x", data)
		})
		stick.On(joystick.LeftY, func(data interface{}) {
			fmt.Println("left_y", data)
		})
		stick.On(joystick.RightX, func(data interface{}) {
			fmt.Println("right_x", data)
		})
		stick.On(joystick.RightY, func(data interface{}) {
			fmt.Println("right_y", data)
		})
	}

	robot := gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{stick},
		work,
	)

	robot.Start()
}
