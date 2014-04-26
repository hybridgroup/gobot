package main

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-opencv"
	"path"
	"runtime"
)

func main() {
	_, currentfile, _, _ := runtime.Caller(0)
	cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")

	opencv := new(gobotOpencv.Opencv)
	opencv.Name = "opencv"

	window := gobotOpencv.NewWindow(opencv)
	window.Name = "window"

	camera := gobotOpencv.NewCamera(opencv)
	camera.Name = "camera"

	work := func() {
		var image *cv.IplImage
		gobot.On(camera.Events["Frame"], func(data interface{}) {
			image = data.(*cv.IplImage)
		})

		go func() {
			for {
				if image != nil {
					i := image.Clone()
					faces := gobotOpencv.DetectFaces(cascade, i)
					i = gobotOpencv.DrawRectangles(i, faces, 0, 255, 0, 5)
					window.ShowImage(i)
				}
			}
		}()
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{opencv},
		Devices:     []gobot.Device{window, camera},
		Work:        work,
	}

	robot.Start()
}
