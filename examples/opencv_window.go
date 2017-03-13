// +build example
//
// Do not build by default.

package main

import (
	cv "github.com/lazywei/go-opencv/opencv"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
)

func main() {
	window := opencv.NewWindowDriver()
	camera := opencv.NewCameraDriver(0)

	work := func() {
		camera.On(camera.Event("frame"), func(data interface{}) {
			window.ShowImage(data.(*cv.IplImage))
		})
	}

	robot := gobot.NewRobot("cameraBot",
		[]gobot.Device{window, camera},
		work,
	)

	robot.Start()
}
