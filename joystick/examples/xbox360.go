package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-joystick"
)

func main() {
	joystickAdaptor := new(gobotJoystick.JoystickAdaptor)
	joystickAdaptor.Name = "xbox360"
	joystickAdaptor.Params = map[string]interface{}{
		"config": "./configs/xbox360_power_a_mini_proex.json",
	}

	joystick := gobotJoystick.NewJoystick(joystickAdaptor)
	joystick.Name = "xbox360"

	work := func() {
		gobot.On(joystick.Events["a_press"], func(data interface{}) {
			fmt.Println("a_press")
		})
		gobot.On(joystick.Events["a_release"], func(data interface{}) {
			fmt.Println("a_release")
		})
		gobot.On(joystick.Events["b_press"], func(data interface{}) {
			fmt.Println("b_press")
		})
		gobot.On(joystick.Events["b_release"], func(data interface{}) {
			fmt.Println("b_release")
		})
		gobot.On(joystick.Events["up"], func(data interface{}) {
			fmt.Println("up", data)
		})
		gobot.On(joystick.Events["down"], func(data interface{}) {
			fmt.Println("down", data)
		})
		gobot.On(joystick.Events["left"], func(data interface{}) {
			fmt.Println("left", data)
		})
		gobot.On(joystick.Events["right"], func(data interface{}) {
			fmt.Println("right", data)
		})
		gobot.On(joystick.Events["left_x"], func(data interface{}) {
			fmt.Println("left_x", data)
		})
		gobot.On(joystick.Events["left_y"], func(data interface{}) {
			fmt.Println("left_y", data)
		})
		gobot.On(joystick.Events["right_x"], func(data interface{}) {
			fmt.Println("right_x", data)
		})
		gobot.On(joystick.Events["right_y"], func(data interface{}) {
			fmt.Println("right_y", data)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{joystickAdaptor},
		Devices:     []gobot.Device{joystick},
		Work:        work,
	}

	robot.Start()
}
