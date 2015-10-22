package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os/exec"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/bebop"
	"github.com/hybridgroup/gobot/platforms/joystick"
)

type pair struct {
	x float64
	y float64
}

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
	gbot := gobot.NewGobot()

	joystickAdaptor := joystick.NewJoystickAdaptor("ps3")
	joystick := joystick.NewJoystickDriver(joystickAdaptor,
		"ps3",
		"./platforms/joystick/configs/dualshock3.json",
	)

	bebopAdaptor := bebop.NewBebopAdaptor("Drone")
	drone := bebop.NewBebopDriver(bebopAdaptor, "Drone")

	work := func() {
		video, _, _ := ffmpeg()

		go func() {
			for {
				if _, err := video.Write(<-drone.Video()); err != nil {
					fmt.Println(err)
          return
				}
			}
		}()

		offset := 32767.0
		rightStick := pair{x: 0, y: 0}
		leftStick := pair{x: 0, y: 0}

		recording := false

		gobot.On(joystick.Event("circle_press"), func(data interface{}) {
			if recording {
				drone.StopRecording()
			} else {
				drone.StartRecording()
			}
			recording = !recording
		})

		gobot.On(joystick.Event("square_press"), func(data interface{}) {
			drone.HullProtection(true)
			drone.TakeOff()
		})
		gobot.On(joystick.Event("triangle_press"), func(data interface{}) {
			drone.Stop()
		})
		gobot.On(joystick.Event("x_press"), func(data interface{}) {
			drone.Land()
		})
		gobot.On(joystick.Event("left_x"), func(data interface{}) {
			val := float64(data.(int16))
			if leftStick.x != val {
				leftStick.x = val
			}
		})
		gobot.On(joystick.Event("left_y"), func(data interface{}) {
			val := float64(data.(int16))
			if leftStick.y != val {
				leftStick.y = val
			}
		})
		gobot.On(joystick.Event("right_x"), func(data interface{}) {
			val := float64(data.(int16))
			if rightStick.x != val {
				rightStick.x = val
			}
		})
		gobot.On(joystick.Event("right_y"), func(data interface{}) {
			val := float64(data.(int16))
			if rightStick.y != val {
				rightStick.y = val
			}
		})

		gobot.Every(10*time.Millisecond, func() {
			pair := leftStick
			if pair.y < -10 {
				drone.Forward(validatePitch(pair.y, offset))
			} else if pair.y > 10 {
				drone.Backward(validatePitch(pair.y, offset))
			} else {
				drone.Forward(0)
			}

			if pair.x > 10 {
				drone.Right(validatePitch(pair.x, offset))
			} else if pair.x < -10 {
				drone.Left(validatePitch(pair.x, offset))
			} else {
				drone.Right(0)
			}
		})

		gobot.Every(10*time.Millisecond, func() {
			pair := rightStick
			if pair.y < -10 {
				drone.Up(validatePitch(pair.y, offset))
			} else if pair.y > 10 {
				drone.Down(validatePitch(pair.y, offset))
			} else {
				drone.Up(0)
			}

			if pair.x > 20 {
				drone.Clockwise(validatePitch(pair.x, offset))
			} else if pair.x < -20 {
				drone.CounterClockwise(validatePitch(pair.x, offset))
			} else {
				drone.Clockwise(0)
			}
		})
	}

	robot := gobot.NewRobot("bebop",
		[]gobot.Connection{joystickAdaptor, bebopAdaptor},
		[]gobot.Device{joystick, drone},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}

func validatePitch(data float64, offset float64) int {
	value := math.Abs(data) / offset
	if value >= 0.1 {
		if value <= 1.0 {
			return int((float64(int(value*100)) / 100) * 100)
		}
		return 100
	}
	return 0
}
