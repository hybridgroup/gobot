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
	stick := joystick.NewDriver(joystickAdaptor, joystick.Xbox360RockBandDrums)

	work := func() {
		_ = stick.On(joystick.RedPress, func(data interface{}) {
			fmt.Println("red_press")
		})
		_ = stick.On(joystick.RedRelease, func(data interface{}) {
			fmt.Println("red_release")
		})
		_ = stick.On(joystick.YellowPress, func(data interface{}) {
			fmt.Println("yellow_press")
		})
		_ = stick.On(joystick.YellowRelease, func(data interface{}) {
			fmt.Println("yellow_release")
		})
		_ = stick.On(joystick.BluePress, func(data interface{}) {
			fmt.Println("blue_press")
		})
		_ = stick.On(joystick.BlueRelease, func(data interface{}) {
			fmt.Println("blue_release")
		})
		_ = stick.On(joystick.GreenPress, func(data interface{}) {
			fmt.Println("green_press")
		})
		_ = stick.On(joystick.GreenRelease, func(data interface{}) {
			fmt.Println("blue_release")
		})
		_ = stick.On(joystick.PedalPress, func(data interface{}) {
			fmt.Println("pedal_press")
		})
		_ = stick.On(joystick.PedalRelease, func(data interface{}) {
			fmt.Println("pedal_release")
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
