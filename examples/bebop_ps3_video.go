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
	"io/ioutil"
	"os/exec"
	"sync/atomic"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/joystick"
	"gobot.io/x/gobot/platforms/parrot/bebop"
)

type pair struct {
	x float64
	y float64
}

var leftX, leftY, rightX, rightY atomic.Value

const offset = 32767.0

func ffmpeg() (stdin io.WriteCloser, stderr io.ReadCloser, err error) {
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "http://localhost:8090/bebop.ffm")

	stderr, err = ffmpeg.StderrPipe()

	if err != nil {
		return
	}

	stdin, err = ffmpeg.StdinPipe()

	if err != nil {
		return
	}

	if err = ffmpeg.Start(); err != nil {
		return
	}

	go func() {
		for {
			buf, err := ioutil.ReadAll(stderr)
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
	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor,
		"./platforms/joystick/configs/dualshock3.json",
	)

	bebopAdaptor := bebop.NewAdaptor()
	drone := bebop.NewDriver(bebopAdaptor)

	work := func() {
		drone.VideoEnable(true)
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

		stick.On(joystick.CirclePress, func(data interface{}) {
			if recording {
				drone.StopRecording()
			} else {
				drone.StartRecording()
			}
			recording = !recording
		})

		stick.On(joystick.SquarePress, func(data interface{}) {
			drone.HullProtection(true)
			drone.TakeOff()
		})
		stick.On(joystick.TrianglePress, func(data interface{}) {
			drone.Stop()
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
				drone.Forward(bebop.ValidatePitch(leftStick.y, offset))
			case leftStick.y > 10:
				drone.Backward(bebop.ValidatePitch(leftStick.y, offset))
			default:
				drone.Forward(0)
			}

			switch {
			case leftStick.x > 10:
				drone.Right(bebop.ValidatePitch(leftStick.x, offset))
			case leftStick.x < -10:
				drone.Left(bebop.ValidatePitch(leftStick.x, offset))
			default:
				drone.Right(0)
			}
		})

		gobot.Every(10*time.Millisecond, func() {
			rightStick := getRightStick()
			switch {
			case rightStick.y < -10:
				drone.Up(bebop.ValidatePitch(rightStick.y, offset))
			case rightStick.y > 10:
				drone.Down(bebop.ValidatePitch(rightStick.y, offset))
			default:
				drone.Up(0)
			}

			switch {
			case rightStick.x > 20:
				drone.Clockwise(bebop.ValidatePitch(rightStick.x, offset))
			case rightStick.x < -20:
				drone.CounterClockwise(bebop.ValidatePitch(rightStick.x, offset))
			default:
				drone.Clockwise(0)
			}
		})
	}

	robot := gobot.NewRobot("bebop",
		[]gobot.Connection{joystickAdaptor, bebopAdaptor},
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
