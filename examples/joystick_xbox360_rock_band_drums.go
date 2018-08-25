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
	stick := joystick.NewDriver(joystickAdaptor, joystick.Xbox360RockBandDrums)

	work := func() {
		stick.On(joystick.RedPress, func(data interface{}) {
			fmt.Println("red_press")
		})
		stick.On(joystick.RedRelease, func(data interface{}) {
			fmt.Println("red_release")
		})
		stick.On(joystick.YellowPress, func(data interface{}) {
			fmt.Println("yellow_press")
		})
		stick.On(joystick.YellowRelease, func(data interface{}) {
			fmt.Println("yellow_release")
		})
		stick.On(joystick.BluePress, func(data interface{}) {
			fmt.Println("blue_press")
		})
		stick.On(joystick.BlueRelease, func(data interface{}) {
			fmt.Println("blue_release")
		})
		stick.On(joystick.GreenPress, func(data interface{}) {
			fmt.Println("green_press")
		})
		stick.On(joystick.GreenRelease, func(data interface{}) {
			fmt.Println("blue_release")
		})
		stick.On(joystick.PedalPress, func(data interface{}) {
			fmt.Println("pedal_press")
		})
		stick.On(joystick.PedalRelease, func(data interface{}) {
			fmt.Println("pedal_release")
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
	}

	robot := gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{stick},
		work,
	)

	robot.Start()
}
