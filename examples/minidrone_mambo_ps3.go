//go:build example
// +build example

//
// Do not build by default.

/*
 How to setup
 You must be using a PS3 or compatible controller, along with
 one of the Parrot Mambo drones to run this example.

 You run the Go program on your computer and communicate
 wirelessly with the Mambo.

 How to run
 Pass the Bluetooth name or address as first param:

	go run examples/minidrone_mambo_ps3.go "Mambo_1234"

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"os"
	"sync/atomic"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/ble"
	"gobot.io/x/gobot/v2/platforms/joystick"
	"gobot.io/x/gobot/v2/platforms/parrot/minidrone"
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

	droneAdaptor := ble.NewClientAdaptor(os.Args[1])
	drone := minidrone.NewDriver(droneAdaptor)

	work := func() {
		leftX.Store(float64(0.0))
		leftY.Store(float64(0.0))
		rightX.Store(float64(0.0))
		rightY.Store(float64(0.0))

		clawOpen := false

		stick.On(joystick.CirclePress, func(data interface{}) {
			if clawOpen {
				drone.ClawControl(0, minidrone.ClawClosed)
				clawOpen = false
			} else {
				drone.ClawControl(0, minidrone.ClawOpen)
				clawOpen = true
			}
		})

		stick.On(joystick.R2Press, func(data interface{}) {
			if clawOpen {
				drone.ClawControl(0, minidrone.ClawClosed)
				clawOpen = false
			} else {
				drone.ClawControl(0, minidrone.ClawOpen)
				clawOpen = true
			}
		})

		stick.On(joystick.TrianglePress, func(data interface{}) {
			drone.HullProtection(true)
			drone.TakeOff()
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
			rightStick := getRightStick()

			switch {
			case rightStick.y < -10:
				drone.Forward(minidrone.ValidatePitch(rightStick.y, offset))
			case rightStick.y > 10:
				drone.Backward(minidrone.ValidatePitch(rightStick.y, offset))
			default:
				drone.Forward(0)
			}

			switch {
			case rightStick.x > 10:
				drone.Right(minidrone.ValidatePitch(rightStick.x, offset))
			case rightStick.x < -10:
				drone.Left(minidrone.ValidatePitch(rightStick.x, offset))
			default:
				drone.Right(0)
			}
		})

		gobot.Every(10*time.Millisecond, func() {
			leftStick := getLeftStick()
			switch {
			case leftStick.y < -10:
				drone.Up(minidrone.ValidatePitch(leftStick.y, offset))
			case leftStick.y > 10:
				drone.Down(minidrone.ValidatePitch(leftStick.y, offset))
			default:
				drone.Up(0)
			}

			switch {
			case leftStick.x > 20:
				drone.Clockwise(minidrone.ValidatePitch(leftStick.x, offset))
			case leftStick.x < -20:
				drone.CounterClockwise(minidrone.ValidatePitch(leftStick.x, offset))
			default:
				drone.Clockwise(0)
			}
		})
	}

	robot := gobot.NewRobot("minidrone",
		[]gobot.Connection{joystickAdaptor, droneAdaptor},
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
