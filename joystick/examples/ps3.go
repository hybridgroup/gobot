package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-joystick"
)

func main() {
	joystickAdaptor := new(gobotJoystick.JoystickAdaptor)
	joystickAdaptor.Name = "ps3"
	joystickAdaptor.Params = map[string]interface{}{
		"config": "./configs/dualshock3.json",
	}

	joystick := gobotJoystick.NewJoystick(joystickAdaptor)
	joystick.Name = "ps3"

	work := func() {
		gobot.On(joystick.Events["square_press"], func(data interface{}) {
			fmt.Println("square_press")
		})
		gobot.On(joystick.Events["square_release"], func(data interface{}) {
			fmt.Println("square_release")
		})
		gobot.On(joystick.Events["triangle_press"], func(data interface{}) {
			fmt.Println("triangle_press")
		})
		gobot.On(joystick.Events["triangle_release"], func(data interface{}) {
			fmt.Println("triangle_release")
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
