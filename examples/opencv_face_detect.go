package main

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/opencv"
	"path"
	"runtime"
)

func main() {
	_, currentfile, _, _ := runtime.Caller(0)
	cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")

	gbot := gobot.NewGobot()

	window := opencv.NewWindowDriver("window")
	camera := opencv.NewCameraDriver("camera", 0)

	work := func() {
		var image *cv.IplImage
		gobot.On(camera.Events["Frame"], func(data interface{}) {
			image = data.(*cv.IplImage)
		})

		go func() {
			for {
				if image != nil {
					i := image.Clone()
					faces := opencv.DetectFaces(cascade, i)
					i = opencv.DrawRectangles(i, faces, 0, 255, 0, 5)
					window.ShowImage(i)
				}
			}
		}()
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("faceBot", []gobot.Connection{}, []gobot.Device{window, camera}, work))

	gbot.Start()
}
