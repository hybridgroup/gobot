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
	stick := joystick.NewDriver(joystickAdaptor,
		"./platforms/joystick/configs/dualshock3.json",
	)

	droneAdaptor := bleclient.NewAdaptor(os.Args[1])
	drone := parrot.NewMinidroneDriver(droneAdaptor)

	work := func() {
		leftX.Store(float64(0.0))
		leftY.Store(float64(0.0))
		rightX.Store(float64(0.0))
		rightY.Store(float64(0.0))

		clawOpen := false

		_ = stick.On(joystick.CirclePress, func(data interface{}) {
			if clawOpen {
				if err := drone.ClawControl(0, parrot.ClawClosed); err != nil {
					fmt.Println(err)
				}
				clawOpen = false
			} else {
				if err := drone.ClawControl(0, parrot.ClawOpen); err != nil {
					fmt.Println(err)
				}
				clawOpen = true
			}
		})

		_ = stick.On(joystick.R2Press, func(data interface{}) {
			if clawOpen {
				if err := drone.ClawControl(0, parrot.ClawClosed); err != nil {
					fmt.Println(err)
				}
				clawOpen = false
			} else {
				if err := drone.ClawControl(0, parrot.ClawOpen); err != nil {
					fmt.Println(err)
				}
				clawOpen = true
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
