// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"math"
	"path"
	"runtime"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
	"gobot.io/x/gobot/platforms/parrot/ardrone"
	"gocv.io/x/gocv"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	_, currentfile, _, _ := runtime.Caller(0)
	cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")
	window := opencv.NewWindowDriver()
	camera := opencv.NewCameraDriver("tcp://192.168.1.1:5555")
	ardroneAdaptor := ardrone.NewAdaptor()
	drone := ardrone.NewDriver(ardroneAdaptor)

	work := func() {
		detect := false
		drone.TakeOff()
		var img gocv.Mat
		camera.On(opencv.Frame, func(data interface{}) {
			img = data.(gocv.Mat)
			if !detect {
				window.IMShow(img)
				window.WaitKey(1)
			}
		})
		drone.On(ardrone.Flying, func(data interface{}) {
			gobot.After(1*time.Second, func() { drone.Up(0.2) })
			gobot.After(2*time.Second, func() { drone.Hover() })
			gobot.After(5*time.Second, func() {
				detect = true
				gobot.Every(300*time.Millisecond, func() {
					drone.Hover()
					i := img
					faces := opencv.DetectObjects(cascade, i)
					biggest := 0
					var face image.Rectangle
					for _, f := range faces {
						if f.Width() > biggest {
							biggest = f.Width()
							face = f
						}
					}
					if face != nil {
						opencv.DrawRectangles(i, []img.Rectangle{face}, 0, 255, 0, 5)
						centerX := float64(img.Size()).X * 0.5
						turn := -(float64(face.Min.X - centerX)) / centerX
						fmt.Println("turning:", turn)
						if turn < 0 {
							drone.Clockwise(math.Abs(turn * 0.4))
						} else {
							drone.CounterClockwise(math.Abs(turn * 0.4))
						}
					}
					window.IMShow(i)
					window.WaitKey(1)
				})
				gobot.After(20*time.Second, func() { drone.Land() })
			})
		})
	}

	robot := gobot.NewRobot("face",
		[]gobot.Connection{ardroneAdaptor},
		[]gobot.Device{window, camera, drone},
		work,
	)

	robot.Start()
}
