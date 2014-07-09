package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
)

func main() {
	gbot := gobot.NewGobot()

	joystickAdaptor := joystick.NewJoystickAdaptor("xbox360")
	joystickDriver := joystick.NewJoystickDriver(joystickAdaptor, "xbox360", "./platforms/joystick/configs/joystick/configs/xbox360_power_a_mini_proex.json")

	work := func() {
		gobot.On(joystickDriver.Event("a_press"), func(data interface{}) {
			fmt.Println("a_press")
		})
		gobot.On(joystickDriver.Event("a_release"), func(data interface{}) {
			fmt.Println("a_release")
		})
		gobot.On(joystickDriver.Event("b_press"), func(data interface{}) {
			fmt.Println("b_press")
		})
		gobot.On(joystickDriver.Event("b_release"), func(data interface{}) {
			fmt.Println("b_release")
		})
		gobot.On(joystickDriver.Event("up"), func(data interface{}) {
			fmt.Println("up", data)
		})
		gobot.On(joystickDriver.Event("down"), func(data interface{}) {
			fmt.Println("down", data)
		})
		gobot.On(joystickDriver.Event("left"), func(data interface{}) {
			fmt.Println("left", data)
		})
		gobot.On(joystickDriver.Event("right"), func(data interface{}) {
			fmt.Println("right", data)
		})
		gobot.On(joystickDriver.Event("left_x"), func(data interface{}) {
			fmt.Println("left_x", data)
		})
		gobot.On(joystickDriver.Event("left_y"), func(data interface{}) {
			fmt.Println("left_y", data)
		})
		gobot.On(joystickDriver.Event("right_x"), func(data interface{}) {
			fmt.Println("right_x", data)
		})
		gobot.On(joystickDriver.Event("right_y"), func(data interface{}) {
			fmt.Println("right_y", data)
		})
	}

	robot := gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{joystickDriver},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
