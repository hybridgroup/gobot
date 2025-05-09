//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/opencv"
	"gocv.io/x/gocv"
)

func main() {
	window := opencv.NewWindowDriver()
	camera := opencv.NewCameraDriver(0)

	work := func() {
		camera.On(opencv.Frame, func(data interface{}) {
			img := data.(gocv.Mat)
			window.ShowImage(img)
			window.WaitKey(1)
		})
	}

	robot := gobot.NewRobot("cameraBot",
		[]gobot.Device{window, camera},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
