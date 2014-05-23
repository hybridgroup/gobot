package main

import (
	"fmt"
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ardrone"
	"github.com/hybridgroup/gobot/platforms/opencv"
	"math"
	"path"
	"runtime"
	"time"
)

func main() {
	gbot := gobot.NewGobot()

	_, currentfile, _, _ := runtime.Caller(0)
	cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")
	window := opencv.NewWindowDriver("window")
	camera := opencv.NewCamera("camera", "tcp://192.168.1.1:5555")
	ardroneAdaptor := ardrone.NewArdroneAdaptor("Drone")
	drone := ardrone.NewArdroneDriver(ardroneAdaptor, "drone")

	work := func() {
		detect := false
		drone.TakeOff()
		var image *cv.IplImage
		gobot.On(camera.Events["Frame"], func(data interface{}) {
			image = data.(*cv.IplImage)
			if detect == false {
				window.ShowImage(image)
			}
		})
		gobot.On(drone.Events["Flying"], func(data interface{}) {
			gobot.After(1*time.Second, func() { drone.Up(0.2) })
			gobot.After(2*time.Second, func() { drone.Hover() })
			gobot.After(5*time.Second, func() {
				detect = true
				gobot.Every(0.3*time.Second, func() {
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
						opencv.DrawRectangles(i, []*cv.Rect{face})
						center_x := float64(image.Width()) * 0.5
						turn := -(float64(face.X()) - center_x) / center_x
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

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("face", []gobot.Connection{ardroneAdaptor}, []gobot.Device{window, camera, drone}, work))

	robot.Start()
}
