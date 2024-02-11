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
	stick := joystick.NewDriver(joystickAdaptor, joystick.Dualshock4)

	work := func() {
		_ = stick.On(joystick.SquarePress, func(data interface{}) {
			fmt.Println("square_press")
		})
		_ = stick.On(joystick.SquareRelease, func(data interface{}) {
			fmt.Println("square_release")
		})
		_ = stick.On(joystick.TrianglePress, func(data interface{}) {
			fmt.Println("triangle_press")
		})
		_ = stick.On(joystick.TriangleRelease, func(data interface{}) {
			fmt.Println("triangle_release")
		})
		_ = stick.On(joystick.CirclePress, func(data interface{}) {
			fmt.Println("circle_press")
		})
		_ = stick.On(joystick.CircleRelease, func(data interface{}) {
			fmt.Println("circle_release")
		})
		_ = stick.On(joystick.XPress, func(data interface{}) {
			fmt.Println("x_press")
		})
		_ = stick.On(joystick.XRelease, func(data interface{}) {
			fmt.Println("x_release")
		})
		_ = stick.On(joystick.HomePress, func(data interface{}) {
			fmt.Println("home_press")
		})
		_ = stick.On(joystick.HomeRelease, func(data interface{}) {
			fmt.Println("home_release")
		})
		_ = stick.On(joystick.SharePress, func(data interface{}) {
			fmt.Println("share_press")
		})
		_ = stick.On(joystick.ShareRelease, func(data interface{}) {
			fmt.Println("share_release")
		})
		_ = stick.On(joystick.OptionsPress, func(data interface{}) {
			fmt.Println("options_press")
		})
		_ = stick.On(joystick.OptionsRelease, func(data interface{}) {
			fmt.Println("options_release")
		})
		_ = stick.On(joystick.L1Press, func(data interface{}) {
			fmt.Println("l1_press")
		})
		_ = stick.On(joystick.L1Release, func(data interface{}) {
			fmt.Println("l1_release")
		})
		_ = stick.On(joystick.L2Press, func(data interface{}) {
			fmt.Println("l2_press")
		})
		_ = stick.On(joystick.L2Release, func(data interface{}) {
			fmt.Println("l2_release")
		})
		_ = stick.On(joystick.R1Press, func(data interface{}) {
			fmt.Println("r1_press")
		})
		_ = stick.On(joystick.R1Release, func(data interface{}) {
			fmt.Println("r1_release")
		})
		_ = stick.On(joystick.R2Press, func(data interface{}) {
			fmt.Println("r2_press")
		})
		_ = stick.On(joystick.R2Release, func(data interface{}) {
			fmt.Println("r2_release")
		})

		_ = stick.On(joystick.UpPress, func(data interface{}) {
			fmt.Println("up_press")
		})
		_ = stick.On(joystick.UpRelease, func(data interface{}) {
			fmt.Println("up_release")
		})
		_ = stick.On(joystick.DownPress, func(data interface{}) {
			fmt.Println("down_press")
		})
		_ = stick.On(joystick.DownRelease, func(data interface{}) {
			fmt.Println("down_release")
		})
		_ = stick.On(joystick.RightPress, func(data interface{}) {
			fmt.Println("right_press")
		})
		_ = stick.On(joystick.RightRelease, func(data interface{}) {
			fmt.Println("right_release")
		})
		_ = stick.On(joystick.LeftPress, func(data interface{}) {
			fmt.Println("left_press")
		})
		_ = stick.On(joystick.LeftRelease, func(data interface{}) {
			fmt.Println("left_release")
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
		_ = stick.On(joystick.L2, func(data interface{}) {
			fmt.Println("L2", data)
		})
		_ = stick.On(joystick.R2, func(data interface{}) {
			fmt.Println("R2", data)
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
