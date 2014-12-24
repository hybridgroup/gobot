package opencv

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*WindowDriver)(nil)

type window interface {
	ShowImage(*cv.IplImage)
}

type WindowDriver struct {
	name   string
	window window
	start  func(*WindowDriver)
}

// NewWindowDriver creates a new window driver with specified name.
// It adds an start function to initialize window
func NewWindowDriver(name string) *WindowDriver {
	return &WindowDriver{
		name: name,
		start: func(w *WindowDriver) {
			w.window = cv.NewWindow(w.Name(), cv.CV_WINDOW_NORMAL)
		},
	}
}

func (w *WindowDriver) Name() string                 { return w.name }
func (w *WindowDriver) Connection() gobot.Connection { return nil }

// Start starts window thread and driver
func (w *WindowDriver) Start() (errs []error) {
	cv.StartWindowThread()
	w.start(w)
	return
}

// Halt returns true if camera is halted successfully
func (w *WindowDriver) Halt() (errs []error) { return }

// ShowImage displays image in window
func (w *WindowDriver) ShowImage(image *cv.IplImage) {
	w.window.ShowImage(image)
}
