package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
)

func main() {
	gbot := gobot.NewGobot()
	joystickAdaptor := joystick.NewJoystickAdaptor("xbox360")
	joystickDriver := joystick.NewJoystickDriver(joystickAdaptor, "xbox360", "../joystick/configs/xbox360_power_a_mini_proex.json")

	work := func() {
		gobot.On(joystickDriver.Events["a_press"], func(data interface{}) {
			fmt.Println("a_press")
		})
		gobot.On(joystickDriver.Events["a_release"], func(data interface{}) {
			fmt.Println("a_release")
		})
		gobot.On(joystickDriver.Events["b_press"], func(data interface{}) {
			fmt.Println("b_press")
		})
		gobot.On(joystickDriver.Events["b_release"], func(data interface{}) {
			fmt.Println("b_release")
		})
		gobot.On(joystickDriver.Events["up"], func(data interface{}) {
			fmt.Println("up", data)
		})
		gobot.On(joystickDriver.Events["down"], func(data interface{}) {
			fmt.Println("down", data)
		})
		gobot.On(joystickDriver.Events["left"], func(data interface{}) {
			fmt.Println("left", data)
		})
		gobot.On(joystickDriver.Events["right"], func(data interface{}) {
			fmt.Println("right", data)
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
