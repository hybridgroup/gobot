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

// NewWindowDriver creates a new window driver with specified name.
// It adds an start function to initialize window
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

// Start starts window thread and driver
func (w *WindowDriver) Start() bool {
	cv.StartWindowThread()
	w.start(w)
	return true
}

// Halt returns true if camera is halted successfully
func (w *WindowDriver) Halt() bool { return true }

// Init returns true if driver is initialized correctly
func (w *WindowDriver) Init() bool { return true }

// ShowImage displays image in window
func (w *WindowDriver) ShowImage(image *cv.IplImage) {
	w.window.ShowImage(image)
}
