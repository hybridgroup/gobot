package opencv

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
)

type WindowDriver struct {
	gobot.Driver
	window *cv.Window
}

func NewWindowDriver(name string) *WindowDriver {
	return &WindowDriver{
		Driver: gobot.Driver{
			Name: name,
		},
	}
}

func (w *WindowDriver) Start() bool {
	cv.StartWindowThread()
	w.window = cv.NewWindow(w.Name, cv.CV_WINDOW_NORMAL)
	return true
}

func (w *WindowDriver) Halt() bool { return true }
func (w *WindowDriver) Init() bool { return true }

func (w *WindowDriver) ShowImage(image *cv.IplImage) {
	w.window.ShowImage(image)
}
