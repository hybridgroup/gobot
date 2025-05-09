//go:build example
// +build example

//
// Do not build by default.

/*
This example will connect to the Parrot Bebop allowing you to control
it using a PS3 controller.

It also streams the drone video to a webpage via ffserver.

This requires you to have both ffmpeg and ffserver installed
on your computer.

In order to run this example you will first need to start ffserver with:

	$ ffserver -f ff.conf

then in a separate terminal run this program:

	$ go run bebop_ps3_video.go

You can view the video feed by navigating to
http://localhost:8090/bebop.mjpeg in a web browser.
*NOTE* firefox works best for viewing the video feed.
*/
package main

import (
	"fmt"
	"io"
	"os/exec"
	"sync/atomic"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/joystick"
	"gobot.io/x/gobot/v2/platforms/parrot/bebop"
)

type pair struct {
	x float64
	y float64
}

var leftX, leftY, rightX, rightY atomic.Value

const offset = 32767.0

func ffmpeg() (io.WriteCloser, io.ReadCloser, error) {
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "http://localhost:8090/bebop.ffm")

	stderr, err := ffmpeg.StderrPipe()
	if err != nil {
		return nil, stderr, err
	}

	stdin, err := ffmpeg.StdinPipe()
	if err != nil {
		return stdin, stderr, err
	}

	if err := ffmpeg.Start(); err != nil {
		return stdin, stderr, err
	}

	go func() {
		for {
			buf, err := io.ReadAll(stderr)
			if err != nil {
				fmt.Println(err)
			}
			if len(buf) > 0 {
				fmt.Println(string(buf))
			}
		}
	}()

	return stdin, stderr, nil
}

func main() {
	joystickAdaptor := joystick.NewAdaptor("0")
	stick := joystick.NewDriver(joystickAdaptor, "dualshock3")

	bebopAdaptor := bebop.NewAdaptor()
	drone := bebop.NewDriver(bebopAdaptor)

	work := func() {
		if err := drone.VideoEnable(true); err != nil {
			fmt.Println(err)
		}
		video, _, _ := ffmpeg()

		go func() {
			for {
				if _, err := video.Write(<-drone.Video()); err != nil {
					fmt.Println(err)
					return
				}
			}
		}()

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
			if err := drone.HullProtection(true); err != nil {
				fmt.Println(err)
			}
			drone.TakeOff()
		})
		_ = stick.On(joystick.TrianglePress, func(data interface{}) {
			if err := drone.Stop(); err != nil {
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
			leftStick := getLeftStick()

			switch {
			case leftStick.y < -10:
				if err := drone.Forward(bebop.ValidatePitch(leftStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			case leftStick.y > 10:
				if err := drone.Backward(bebop.ValidatePitch(leftStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Forward(0); err != nil {
					fmt.Println(err)
				}
			}

			switch {
			case leftStick.x > 10:
				if err := drone.Right(bebop.ValidatePitch(leftStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			case leftStick.x < -10:
				if err := drone.Left(bebop.ValidatePitch(leftStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Right(0); err != nil {
					fmt.Println(err)
				}
			}
		})

		gobot.Every(10*time.Millisecond, func() {
			rightStick := getRightStick()
			switch {
			case rightStick.y < -10:
				if err := drone.Up(bebop.ValidatePitch(rightStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			case rightStick.y > 10:
				if err := drone.Down(bebop.ValidatePitch(rightStick.y, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Up(0); err != nil {
					fmt.Println(err)
				}
			}

			switch {
			case rightStick.x > 20:
				if err := drone.Clockwise(bebop.ValidatePitch(rightStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			case rightStick.x < -20:
				if err := drone.CounterClockwise(bebop.ValidatePitch(rightStick.x, offset)); err != nil {
					fmt.Println(err)
				}
			default:
				if err := drone.Clockwise(0); err != nil {
					fmt.Println(err)
				}
			}
		})
	}

	robot := gobot.NewRobot("bebop",
		[]gobot.Connection{joystickAdaptor, bebopAdaptor},
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
