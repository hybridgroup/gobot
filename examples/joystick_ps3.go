package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/joystick"
)

func main() {
	joystickAdaptor := joystick.NewJoystickAdaptor()
	joystickAdaptor.Name = "ps3"
	joystickAdaptor.Params = map[string]interface{}{
		"config": "../joystick/configs/dualshock3.json",
	}

	joystickDriver := joystick.NewJoystickDriver(joystickAdaptor)
	joystickDriver.Name = "ps3"

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

	robot := gobot.Robot{
		Connections: []gobot.Connection{joystickAdaptor},
		Devices:     []gobot.Device{joystickDriver},
		Work:        work,
	}

	robot.Start()
}
