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

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gobot.io/x/gobot/platforms/joystick"
)

type pair struct {
	x float64
	y float64
}

var leftX, leftY, rightX, rightY atomic.Value

const offset = 32767.0

func main() {
	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor, "dualshock3")

	drone := tello.NewDriver("8888")

	work := func() {
		leftX.Store(float64(0.0))
		leftY.Store(float64(0.0))
		rightX.Store(float64(0.0))
		rightY.Store(float64(0.0))

		stick.On(joystick.TrianglePress, func(data interface{}) {
			drone.TakeOff()
		})

		stick.On(joystick.XPress, func(data interface{}) {
			drone.Land()
		})

		stick.On(joystick.UpPress, func(data interface{}) {
			fmt.Println("FrontFlip")
			drone.FrontFlip()
		})

		stick.On(joystick.DownPress, func(data interface{}) {
			fmt.Println("BackFlip")
			drone.BackFlip()
		})

		stick.On(joystick.RightPress, func(data interface{}) {
			fmt.Println("RightFlip")
			drone.RightFlip()
		})

		stick.On(joystick.LeftPress, func(data interface{}) {
			fmt.Println("LeftFlip")
			drone.LeftFlip()
		})

		stick.On(joystick.LeftX, func(data interface{}) {
			val := float64(data.(int16))
			leftX.Store(val)
		})

		stick.On(joystick.LeftY, func(data interface{}) {
			val := float64(data.(int16))
			leftY.Store(val)
		})

		stick.On(joystick.RightX, func(data interface{}) {
			val := float64(data.(int16))
			rightX.Store(val)
		})

		stick.On(joystick.RightY, func(data interface{}) {
			val := float64(data.(int16))
			rightY.Store(val)
		})

		gobot.Every(50*time.Millisecond, func() {
			rightStick := getRightStick()

			switch {
			case rightStick.y < -10:
				drone.Forward(tello.ValidatePitch(rightStick.y, offset))
			case rightStick.y > 10:
				drone.Backward(tello.ValidatePitch(rightStick.y, offset))
			default:
				drone.Forward(0)
			}

			switch {
			case rightStick.x > 10:
				drone.Right(tello.ValidatePitch(rightStick.x, offset))
			case rightStick.x < -10:
				drone.Left(tello.ValidatePitch(rightStick.x, offset))
			default:
				drone.Right(0)
			}
		})

		gobot.Every(50*time.Millisecond, func() {
			leftStick := getLeftStick()
			switch {
			case leftStick.y < -10:
				drone.Up(tello.ValidatePitch(leftStick.y, offset))
			case leftStick.y > 10:
				drone.Down(tello.ValidatePitch(leftStick.y, offset))
			default:
				drone.Up(0)
			}

			switch {
			case leftStick.x > 20:
				drone.Clockwise(tello.ValidatePitch(leftStick.x, offset))
			case leftStick.x < -20:
				drone.CounterClockwise(tello.ValidatePitch(leftStick.x, offset))
			default:
				drone.Clockwise(0)
			}
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{stick, drone},
		work,
	)

	robot.Start()
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
