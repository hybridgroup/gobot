//go:build gocv
// +build gocv

package opencv

import (
	"gobot.io/x/gobot/v2"
	"gocv.io/x/gocv"
)

type window interface {
	IMShow(gocv.Mat)
}

// WindowDriver is the Gobot Driver for the OpenCV window
type WindowDriver struct {
	name   string
	window window
	start  func(*WindowDriver)
}

// NewWindowDriver creates a new window driver.
// It adds an start function to initialize window
func NewWindowDriver() *WindowDriver {
	return &WindowDriver{
		name: "Window",
		start: func(w *WindowDriver) {
			w.window = gocv.NewWindow(w.Name())
		},
	}
}

// Name returns the Driver name
func (w *WindowDriver) Name() string { return w.name }

// SetName sets the Driver name
func (w *WindowDriver) SetName(n string) { w.name = n }

// Connection returns the Driver's connection
func (w *WindowDriver) Connection() gobot.Connection { return nil }

// Start starts window thread and driver
func (w *WindowDriver) Start() error {
	w.start(w)
	return nil
}

// Halt returns true if camera is halted successfully
func (w *WindowDriver) Halt() error { return nil }

// ShowImage displays image in window
func (w *WindowDriver) ShowImage(img gocv.Mat) {
	w.window.IMShow(img)
}

// WaitKey gives window a chance to redraw, and checks for keyboard input
func (w *WindowDriver) WaitKey(pause int) int {
	return gocv.WaitKey(pause)
}
