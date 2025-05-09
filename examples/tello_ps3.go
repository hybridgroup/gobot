//go:build example
// +build example

//
// Do not build by default.

/*
How to setup
You must be using a PS3 or compatible controller, along with a DJI Tello drone to run this example.

You run the Go program on your computer and communicate wirelessly via WiFi with the Tello.

How to run

	go run examples/tello_ps3.go
*/

package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/dji/tello"
	"gobot.io/x/gobot/v2/platforms/joystick"
)

type pair struct {
	x float64
	y float64
}

var leftX, leftY, rightX, rightY atomic.Value

const offset = 32767.0

func main() {
	joystickAdaptor := joystick.NewAdaptor("0")
	stick := joystick.NewDriver(joystickAdaptor, "dualshock3")

	drone := tello.NewDriver("8888")

	work := func() {
		leftX.Store(float64(0.0))
		leftY.Store(float64(0.0))
		rightX.Store(float64(0.0))
		rightY.Store(float64(0.0))

		_ = stick.On(joystick.TrianglePress, func(data interface{}) {
			if err := drone.TakeOff(); err != nil {
				fmt.Println(err)
			}
		})

		_ = stick.On(joystick.XPress, func(data interface{}) {
			if err := drone.Land(); err != nil {
				fmt.Println(err)
			}
		})

		_ = stick.On(joystick.UpPress, func(data interface{}) {
			fmt.Println("FrontFlip")
			if err := drone.FrontFlip(); err != nil {
				fmt.Println(err)
			}
		})

		_ = stick.On(joystick.DownPress, func(data interface{}) {
			fmt.Println("BackFlip")
			if err := drone.BackFlip(); err != nil {
				fmt.Println(err)
			}
		})

		_ = stick.On(joystick.RightPress, func(data interface{}) {
			fmt.Println("RightFlip")
			if err := drone.RightFlip(); err != nil {
				fmt.Println(err)
			}
		})

		_ = stick.On(joystick.LeftPress, func(data interface{}) {
			fmt.Println("LeftFlip")
			if err := drone.LeftFlip(); err != nil {
				fmt.Println(err)
			}
		})

		_ = stick.On(joystick.LeftX, func(data interface{}) {
			val := float64(data.(int))
			leftX.Store(val)
		})

		_ = stick.On(joystick.LeftY, func(data interface{}) {
			val := float64(data.(int))
			leftY.Store(val)
		})

		_ = stick.On(joystick.RightX, func(data interface{}) {
			val := float64(data.(int))
			rightX.Store(val)
		})

		_ = stick.On(joystick.RightY, func(data interface{}) {
			val := float64(data.(int))
			rightY.Store(val)
		})

		gobot.Every(50*time.Millisecond, func() {
			rightStick := getRightStick()

			switch {
			case rightStick.y < -10:
				if err := drone.Forward(tello.ValidatePitch(rightStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			case rightStick.y > 10:
				if err := drone.Backward(tello.ValidatePitch(rightStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Forward(0); err != nil {
					fmt.Println(err)
				}
			}

			switch {
			case rightStick.x > 10:
				if err := drone.Right(tello.ValidatePitch(rightStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			case rightStick.x < -10:
				if err := drone.Left(tello.ValidatePitch(rightStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Right(0); err != nil {
					fmt.Println(err)
				}
			}
		})

		gobot.Every(50*time.Millisecond, func() {
			leftStick := getLeftStick()
			switch {
			case leftStick.y < -10:
				if err := drone.Up(tello.ValidatePitch(leftStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			case leftStick.y > 10:
				if err := drone.Down(tello.ValidatePitch(leftStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Up(0); err != nil {
					fmt.Println(err)
				}
			}

			switch {
			case leftStick.x > 20:
				if err := drone.Clockwise(tello.ValidatePitch(leftStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			case leftStick.x < -20:
				if err := drone.CounterClockwise(tello.ValidatePitch(leftStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Clockwise(0); err != nil {
					fmt.Println(err)
				}
			}
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{stick, drone},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}

func getLeftStick() pair {
	s := pair{x: 0, y: 0}
	s.x = leftX.Load().(float64)
	s.y = leftY.Load().(float64)
	return s
}

func getRightStick() pair {
	s := pair{x: 0, y: 0}
	s.x = rightX.Load().(float64)
	s.y = rightY.Load().(float64)
	return s
}
