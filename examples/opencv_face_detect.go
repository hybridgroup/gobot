// +build example
//
// Do not build by default.

package main

import (
	"path"
	"runtime"
	"time"

	cv "github.com/lazywei/go-opencv/opencv"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
)

func main() {
	_, currentfile, _, _ := runtime.Caller(0)
	cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")

	window := opencv.NewWindowDriver()
	camera := opencv.NewCameraDriver(0)

	work := func() {
		var image *cv.IplImage

		camera.On(opencv.Frame, func(data interface{}) {
			image = data.(*cv.IplImage)
		})

		gobot.Every(500*time.Millisecond, func() {
			if image != nil {
				i := image.Clone()
				faces := opencv.DetectFaces(cascade, i)
				i = opencv.DrawRectangles(i, faces, 0, 255, 0, 5)
				window.ShowImage(i)
			}

		})
	}

	robot := gobot.NewRobot("faceBot",
		[]gobot.Connection{},
		[]gobot.Device{window, camera},
		work,
	)

	robot.Start()
}
