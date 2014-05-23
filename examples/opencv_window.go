package main

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/opencv"
)

func main() {
	gbot := gobot.NewGobot()

	window := opencv.NewWindowDriver("window")
	camera := opencv.NewCameraDriver("camera", 0)

	work := func() {
		gobot.On(camera.Events["Frame"], func(data interface{}) {
			window.ShowImage(data.(*cv.IplImage))
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("cameraBot", []gobot.Connection{}, []gobot.Device{window, camera}, work))

	gbot.Start()}
}
