package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
)

func main() {
	gbot := gobot.NewGobot()

	joystickAdaptor := joystick.NewJoystickAdaptor("xbox360")
	stick := joystick.NewJoystickDriver(joystickAdaptor,
		"xbox360",
		"./platforms/joystick/configs/joystick/configs/xbox360_power_a_mini_proex.json",
	)

	work := func() {
		stick.On(joystick.APress, func(data interface{}) {
			fmt.Println("a_press")
		})
		stick.On(joystick.ARelease, func(data interface{}) {
			fmt.Println("a_release")
		})
		stick.On(joystick.BPress, func(data interface{}) {
			fmt.Println("b_press")
		})
		stick.On(joystick.BRelease, func(data interface{}) {
			fmt.Println("b_release")
		})
		stick.On(joystick.Up, func(data interface{}) {
			fmt.Println("up", data)
		})
		stick.On(joystick.Down, func(data interface{}) {
			fmt.Println("down", data)
		})
		stick.On(joystick.Left, func(data interface{}) {
			fmt.Println("left", data)
		})
		stick.On(joystick.Right, func(data interface{}) {
			fmt.Println("right", data)
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

	gbot.AddRobot(robot)

	gbot.Start()
}
