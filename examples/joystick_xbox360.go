//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/joystick"
)

func main() {
	joystickAdaptor := joystick.NewAdaptor("0")
	stick := joystick.NewDriver(joystickAdaptor, joystick.Xbox360)

	work := func() {
		_ = stick.On(joystick.APress, func(data interface{}) {
			fmt.Println("a_press")
		})
		_ = stick.On(joystick.ARelease, func(data interface{}) {
			fmt.Println("a_release")
		})
		_ = stick.On(joystick.BPress, func(data interface{}) {
			fmt.Println("b_press")
		})
		_ = stick.On(joystick.BRelease, func(data interface{}) {
			fmt.Println("b_release")
		})
		_ = stick.On(joystick.UpPress, func(data interface{}) {
			fmt.Println("up", data)
		})
		_ = stick.On(joystick.DownPress, func(data interface{}) {
			fmt.Println("down", data)
		})
		_ = stick.On(joystick.LeftPress, func(data interface{}) {
			fmt.Println("left", data)
		})
		_ = stick.On(joystick.RightPress, func(data interface{}) {
			fmt.Println("right", data)
		})
		_ = stick.On(joystick.LeftX, func(data interface{}) {
			fmt.Println("left_x", data)
		})
		_ = stick.On(joystick.LeftY, func(data interface{}) {
			fmt.Println("left_y", data)
		})
		_ = stick.On(joystick.RightX, func(data interface{}) {
			fmt.Println("right_x", data)
		})
		_ = stick.On(joystick.RightY, func(data interface{}) {
			fmt.Println("right_y", data)
		})
	}

	robot := gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{stick},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
