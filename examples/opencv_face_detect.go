// +build example
//
// Do not build by default.

package main

import (
	"path"
	"runtime"
	//"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
	"github.com/hybridgroup/gocv"
)

func main() {
	_, currentfile, _, _ := runtime.Caller(0)
	cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")

	window := opencv.NewWindowDriver()
	camera := opencv.NewCameraDriver(0)

	work := func() {
		camera.On(opencv.Frame, func(data interface{}) {
			i := data.(gocv.Mat)
			faces := opencv.DetectFaces(cascade, i)
			opencv.DrawRectangles(i, faces, 0, 255, 0, 5)
			window.ShowImage(i)
			window.WaitKey(1)	
		})
	}

	robot := gobot.NewRobot("faceBot",
		[]gobot.Connection{},
		[]gobot.Device{window, camera},
		work,
	)

	robot.Start()
}
