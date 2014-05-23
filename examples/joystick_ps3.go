package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
)

func main() {
	gbot := gobot.NewGobot()
	joystickAdaptor := joystick.NewJoystickAdaptor("ps3")
	joystickDriver := joystick.NewJoystickDriver(joystickAdaptor, "ps3", "../joystick/configs/dualshock3.json")

	work := func() {
		gobot.On(joystickDriver.Events["square_press"], func(data interface{}) {
			fmt.Println("square_press")
		})
		gobot.On(joystickDriver.Events["square_release"], func(data interface{}) {
			fmt.Println("square_release")
		})
		gobot.On(joystickDriver.Events["triangle_press"], func(data interface{}) {
			fmt.Println("triangle_press")
		})
		gobot.On(joystickDriver.Events["triangle_release"], func(data interface{}) {
			fmt.Println("triangle_release")
		})
		gobot.On(joystickDriver.Events["left_x"], func(data interface{}) {
			fmt.Println("left_x", data)
		})
		gobot.On(joystickDriver.Events["left_y"], func(data interface{}) {
			fmt.Println("left_y", data)
		})
		gobot.On(joystickDriver.Events["right_x"], func(data interface{}) {
			fmt.Println("right_x", data)
		})
		gobot.On(joystickDriver.Events["right_y"], func(data interface{}) {
			fmt.Println("right_y", data)
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("joystickBot", []gobot.Connection{joystickAdaptor}, []gobot.Device{joystickDriver}, work))

	gbot.Start()
}
