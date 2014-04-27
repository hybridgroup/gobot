package main

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/opencv"
)

func main() {

	window := opencv.NewWindowDriver()
	window.Name = "window"

	camera := opencv.NewCameraDriver()
	camera.Name = "camera"

	work := func() {
		gobot.On(camera.Events["Frame"], func(data interface{}) {
			window.ShowImage(data.(*cv.IplImage))
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{},
		Devices:     []gobot.Device{window, camera},
		Work:        work,
	}

	robot.Start()
}
