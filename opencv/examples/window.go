package main

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-opencv"
)

func main() {

	opencv := new(gobotOpencv.Opencv)
	opencv.Name = "opencv"

	window := gobotOpencv.NewWindow(opencv)
	window.Name = "window"

	camera := gobotOpencv.NewCamera(opencv)
	camera.Name = "camera"

	work := func() {
		gobot.On(camera.Events["Frame"], func(data interface{}) {
			window.ShowImage(data.(*cv.IplImage))
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{opencv},
		Devices:     []gobot.Device{window, camera},
		Work:        work,
	}

	robot.Start()
}
