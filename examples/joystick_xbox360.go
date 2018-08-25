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
	stick := joystick.NewDriver(joystickAdaptor, joystick.Xbox360)

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
		stick.On(joystick.UpPress, func(data interface{}) {
			fmt.Println("up", data)
		})
		stick.On(joystick.DownPress, func(data interface{}) {
			fmt.Println("down", data)
		})
		stick.On(joystick.LeftPress, func(data interface{}) {
			fmt.Println("left", data)
		})
		stick.On(joystick.RightPress, func(data interface{}) {
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

	robot.Start()
}
