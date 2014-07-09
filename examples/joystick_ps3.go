package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
)

func main() {
	gbot := gobot.NewGobot()

	joystickAdaptor := joystick.NewJoystickAdaptor("ps3")
	joystickDriver := joystick.NewJoystickDriver(joystickAdaptor, "ps3", "./platforms/joystick/configs/dualshock3.json")

	work := func() {
		gobot.On(joystickDriver.Event("square_press"), func(data interface{}) {
			fmt.Println("square_press")
		})
		gobot.On(joystickDriver.Event("square_release"), func(data interface{}) {
			fmt.Println("square_release")
		})
		gobot.On(joystickDriver.Event("triangle_press"), func(data interface{}) {
			fmt.Println("triangle_press")
		})
		gobot.On(joystickDriver.Event("triangle_release"), func(data interface{}) {
			fmt.Println("triangle_release")
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
