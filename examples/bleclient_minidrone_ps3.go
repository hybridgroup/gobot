//go:build example
// +build example

//
// Do not build by default.

/*
 How to setup
 You must be using a PS3 or compatible controller, along with
 any of the Parrot Minidrone drones to run this example.

 You run the Go program on your computer and communicate
 wirelessly with the Parrot Minidrone.

 How to run
 Pass the Bluetooth name or address as first param:

	go run examples/minidrone_ps3.go "Travis_1234"

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/parrot"
	"gobot.io/x/gobot/v2/platforms/bleclient"
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

	droneAdaptor := bleclient.NewAdaptor(os.Args[1])
	drone := parrot.NewMinidroneDriver(droneAdaptor)

	work := func() {
		leftX.Store(float64(0.0))
		leftY.Store(float64(0.0))
		rightX.Store(float64(0.0))
		rightY.Store(float64(0.0))

		recording := false

		_ = stick.On(joystick.CirclePress, func(data interface{}) {
			if recording {
				if err := drone.StopRecording(); err != nil {
					fmt.Println(err)
				}
			} else {
				if err := drone.StartRecording(); err != nil {
					fmt.Println(err)
				}
			}
			recording = !recording
		})

		_ = stick.On(joystick.SquarePress, func(data interface{}) {
			if err := drone.Stop(); err != nil {
				fmt.Println(err)
			}
		})

		_ = stick.On(joystick.TrianglePress, func(data interface{}) {
			if err := drone.HullProtection(true); err != nil {
				fmt.Println(err)
			}
			if err := drone.TakeOff(); err != nil {
				fmt.Println(err)
			}
		})

		_ = stick.On(joystick.XPress, func(data interface{}) {
			if err := drone.Land(); err != nil {
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

		gobot.Every(10*time.Millisecond, func() {
			rightStick := getRightStick()

			switch {
			case rightStick.y < -10:
				if err := drone.Forward(parrot.ValidatePitch(rightStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			case rightStick.y > 10:
				if err := drone.Backward(parrot.ValidatePitch(rightStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Forward(0); err != nil {
					fmt.Println(err)
				}
			}

			switch {
			case rightStick.x > 10:
				if err := drone.Right(parrot.ValidatePitch(rightStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			case rightStick.x < -10:
				if err := drone.Left(parrot.ValidatePitch(rightStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Right(0); err != nil {
					fmt.Println(err)
				}
			}
		})

		gobot.Every(10*time.Millisecond, func() {
			leftStick := getLeftStick()
			switch {
			case leftStick.y < -10:
				if err := drone.Up(parrot.ValidatePitch(leftStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			case leftStick.y > 10:
				if err := drone.Down(parrot.ValidatePitch(leftStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Up(0); err != nil {
					fmt.Println(err)
				}
			}

			switch {
			case leftStick.x > 20:
				if err := drone.Clockwise(parrot.ValidatePitch(leftStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			case leftStick.x < -20:
				if err := drone.CounterClockwise(parrot.ValidatePitch(leftStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Clockwise(0); err != nil {
					fmt.Println(err)
				}
			}
		})
	}

	robot := gobot.NewRobot("minidrone",
		[]gobot.Connection{joystickAdaptor, droneAdaptor},
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
