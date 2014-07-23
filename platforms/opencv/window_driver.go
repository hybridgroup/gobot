package opencv

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
)

type WindowDriver struct {
	gobot.Driver
	window window
	start  func(*WindowDriver)
}

func NewWindowDriver(name string) *WindowDriver {
	return &WindowDriver{
		Driver: *gobot.NewDriver(
			name,
			"WindowDriver",
		),
		start: func(w *WindowDriver) {
			w.window = cv.NewWindow(w.Name(), cv.CV_WINDOW_NORMAL)
		},
	}
}

func (w *WindowDriver) Start() bool {
	cv.StartWindowThread()
	w.start(w)
	return true
}

func (w *WindowDriver) Halt() bool { return true }
func (w *WindowDriver) Init() bool { return true }

func (w *WindowDriver) ShowImage(image *cv.IplImage) {
	w.window.ShowImage(image)
}
