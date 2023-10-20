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
	joystick := joystick.NewDriver(joystickAdaptor, joystick.XboxOne)

	work := func() {
		// start button
		joystick.On(joystick.Event("start_press"), func(data interface{}) {
			fmt.Println("start_press")
		})
		joystick.On(joystick.Event("start_release"), func(data interface{}) {
			fmt.Println("start_release")
		})

		// back button
		joystick.On(joystick.Event("back_press"), func(data interface{}) {
			fmt.Println("back_press")
		})
		joystick.On(joystick.Event("back_release"), func(data interface{}) {
			fmt.Println("back_release")
		})

		// a button
		joystick.On(joystick.Event("a_press"), func(data interface{}) {
			fmt.Println("a_press")
		})
		joystick.On(joystick.Event("a_release"), func(data interface{}) {
			fmt.Println("a_release")
		})

		// b button
		joystick.On(joystick.Event("b_press"), func(data interface{}) {
			fmt.Println("b_press")
		})
		joystick.On(joystick.Event("b_release"), func(data interface{}) {
			fmt.Println("b_release")
		})

		// x button
		joystick.On(joystick.Event("x_press"), func(data interface{}) {
			fmt.Println("x_press")
		})
		joystick.On(joystick.Event("x_release"), func(data interface{}) {
			fmt.Println("x_release")
		})

		// y button
		joystick.On(joystick.Event("y_press"), func(data interface{}) {
			fmt.Println("y_press")
		})
		joystick.On(joystick.Event("y_release"), func(data interface{}) {
			fmt.Println("y_release")
		})

		// up dpad
		joystick.On(joystick.Event("up_press"), func(data interface{}) {
			fmt.Println("up_press", data)
		})
		joystick.On(joystick.Event("up_release"), func(data interface{}) {
			fmt.Println("up_release", data)
		})

		// down dpad
		joystick.On(joystick.Event("down_press"), func(data interface{}) {
			fmt.Println("down_press")
		})
		joystick.On(joystick.Event("down_release"), func(data interface{}) {
			fmt.Println("down_release")
		})

		// left dpad
		joystick.On(joystick.Event("left_press"), func(data interface{}) {
			fmt.Println("left_press")
		})
		joystick.On(joystick.Event("left_release"), func(data interface{}) {
			fmt.Println("left_release")
		})

		// right dpad
		joystick.On(joystick.Event("right_press"), func(data interface{}) {
			fmt.Println("right_press")
		})
		joystick.On(joystick.Event("right_release"), func(data interface{}) {
			fmt.Println("right_release")
		})

		// rt trigger
		joystick.On(joystick.Event("rt"), func(data interface{}) {
			fmt.Println("rt", data)
		})

		// lt trigger
		joystick.On(joystick.Event("lt"), func(data interface{}) {
			fmt.Println("lt", data)
		})

		// lb button
		joystick.On(joystick.Event("lb_press"), func(data interface{}) {
			fmt.Println("lb_press")
		})
		joystick.On(joystick.Event("lb_release"), func(data interface{}) {
			fmt.Println("lb_release")
		})

		// rb button
		joystick.On(joystick.Event("rb_press"), func(data interface{}) {
			fmt.Println("rb_press")
		})
		joystick.On(joystick.Event("rb_release"), func(data interface{}) {
			fmt.Println("rb_release")
		})

		// rx stick
		joystick.On(joystick.Event("right_x"), func(data interface{}) {
			fmt.Println("right_x", data)
		})

		// ry stick
		joystick.On(joystick.Event("right_y"), func(data interface{}) {
			fmt.Println("right_y", data)
		})

		// right_stick button
		joystick.On(joystick.Event("right_stick_press"), func(data interface{}) {
			fmt.Println("right_stick_press")
		})
		joystick.On(joystick.Event("right_stick_release"), func(data interface{}) {
			fmt.Println("right_stick_release")
		})

		// lx stick
		joystick.On(joystick.Event("left_x"), func(data interface{}) {
			fmt.Println("left_x", data)
		})

		// ly stick
		joystick.On(joystick.Event("left_y"), func(data interface{}) {
			fmt.Println("left_y", data)
		})

		// left_stick button
		joystick.On(joystick.Event("left_stick_press"), func(data interface{}) {
			fmt.Println("left_stick_press")
		})
		joystick.On(joystick.Event("left_stick_release"), func(data interface{}) {
			fmt.Println("left_stick_release")
		})
	}

	robot := gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{joystick},
		work,
	)

	robot.Start()
}
