// +build example
//
// Do not build by default.

package main

import (
	"sync/atomic"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/joystick"
	"gobot.io/x/gobot/platforms/parrot/ardrone"
)

type pair struct {
	x float64
	y float64
}

var leftX, leftY, rightX, rightY atomic.Value

const offset = 32767.0

func main() {
	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor,
		"./platforms/joystick/configs/dualshock3.json",
	)

	ardroneAdaptor := ardrone.NewAdaptor()
	drone := ardrone.NewDriver(ardroneAdaptor)

	leftX.Store(float64(0.0))
	leftY.Store(float64(0.0))
	rightX.Store(float64(0.0))
	rightY.Store(float64(0.0))

	work := func() {
		stick.On(joystick.SquarePress, func(data interface{}) {
			drone.TakeOff()
		})

		stick.On(joystick.TrianglePress, func(data interface{}) {
			drone.Hover()
		})

		stick.On(joystick.XPress, func(data interface{}) {
			drone.Land()
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

		gobot.Every(10*time.Millisecond, func() {
			leftStick := getLeftStick()

			switch {
			case leftStick.y < -10:
				drone.Forward(ardrone.ValidatePitch(leftStick.y, offset))
			case leftStick.y > 10:
				drone.Backward(ardrone.ValidatePitch(leftStick.y, offset))
			default:
				drone.Forward(0)
			}

			switch {
			case leftStick.x > 10:
				drone.Right(ardrone.ValidatePitch(leftStick.x, offset))
			case leftStick.x < -10:
				drone.Left(ardrone.ValidatePitch(leftStick.x, offset))
			default:
				drone.Right(0)
			}
		})

		gobot.Every(10*time.Millisecond, func() {
			rightStick := getRightStick()

			switch {
			case rightStick.y < -10:
				drone.Up(ardrone.ValidatePitch(rightStick.y, offset))
			case rightStick.y > 10:
				drone.Down(ardrone.ValidatePitch(rightStick.y, offset))
			default:
				drone.Up(0)
			}

			switch {
			case rightStick.x > 20:
				drone.Clockwise(ardrone.ValidatePitch(rightStick.x, offset))
			case rightStick.x < -20:
				drone.CounterClockwise(ardrone.ValidatePitch(rightStick.x, offset))
			default:
				drone.Clockwise(0)
			}
		})
	}

	robot := gobot.NewRobot("ardrone",
		[]gobot.Connection{joystickAdaptor, ardroneAdaptor},
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
