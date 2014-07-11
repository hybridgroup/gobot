package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
)

func main() {
	gbot := gobot.NewGobot()

	joystickAdaptor := joystick.NewJoystickAdaptor("xbox360")
	joystick := joystick.NewJoystickDriver(joystickAdaptor,
		"xbox360",
		"./platforms/joystick/configs/joystick/configs/xbox360_power_a_mini_proex.json",
	)

	work := func() {
		gobot.On(joystick.Event("a_press"), func(data interface{}) {
			fmt.Println("a_press")
		})
		gobot.On(joystick.Event("a_release"), func(data interface{}) {
			fmt.Println("a_release")
		})
		gobot.On(joystick.Event("b_press"), func(data interface{}) {
			fmt.Println("b_press")
		})
		gobot.On(joystick.Event("b_release"), func(data interface{}) {
			fmt.Println("b_release")
		})
		gobot.On(joystick.Event("up"), func(data interface{}) {
			fmt.Println("up", data)
		})
		gobot.On(joystick.Event("down"), func(data interface{}) {
			fmt.Println("down", data)
		})
		gobot.On(joystick.Event("left"), func(data interface{}) {
			fmt.Println("left", data)
		})
		gobot.On(joystick.Event("right"), func(data interface{}) {
			fmt.Println("right", data)
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
