package main

import (
	"fmt"
	"math"
	"path"
	"runtime"
	"time"

	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ardrone"
	"github.com/hybridgroup/gobot/platforms/opencv"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	gbot := gobot.NewGobot()

	_, currentfile, _, _ := runtime.Caller(0)
	cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")
	window := opencv.NewWindowDriver("window")
	camera := opencv.NewCameraDriver("camera", "tcp://192.168.1.1:5555")
	ardroneAdaptor := ardrone.NewArdroneAdaptor("Drone")
	drone := ardrone.NewArdroneDriver(ardroneAdaptor, "drone")

	work := func() {
		detect := false
		drone.TakeOff()
		var image *cv.IplImage
		gobot.On(camera.Event("frame"), func(data interface{}) {
			image = data.(*cv.IplImage)
			if detect == false {
				window.ShowImage(image)
			}
		})
		gobot.On(drone.Event("flying"), func(data interface{}) {
			gobot.After(1*time.Second, func() { drone.Up(0.2) })
			gobot.After(2*time.Second, func() { drone.Hover() })
			gobot.After(5*time.Second, func() {
				detect = true
				gobot.Every(300*time.Millisecond, func() {
					drone.Hover()
					i := image
					faces := opencv.DetectFaces(cascade, i)
					biggest := 0
					var face *cv.Rect
					for _, f := range faces {
						if f.Width() > biggest {
							biggest = f.Width()
							face = f
						}
					}
					if face != nil {
						opencv.DrawRectangles(i, []*cv.Rect{face}, 0, 255, 0, 5)
						centerX := float64(image.Width()) * 0.5
						turn := -(float64(face.X()) - centerX) / centerX
						fmt.Println("turning:", turn)
						if turn < 0 {
							drone.Clockwise(math.Abs(turn * 0.4))
						} else {
							drone.CounterClockwise(math.Abs(turn * 0.4))
						}
					}
					window.ShowImage(i)
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

	gbot.AddRobot(robot)

	gbot.Start()
}
