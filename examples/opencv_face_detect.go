package main

import (
	"path"
	"runtime"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/opencv"
	cv "github.com/lazywei/go-opencv/opencv"
)

func main() {
	_, currentfile, _, _ := runtime.Caller(0)
	cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")

	gbot := gobot.NewGobot()

	window := opencv.NewWindowDriver("window")
	camera := opencv.NewCameraDriver("camera", 0)

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

	gbot.AddRobot(robot)

	gbot.Start()
}
