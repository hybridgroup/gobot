package opencv

import (
	"gocv.io/x/gocv"
)

type testCapture struct{}

func (c *testCapture) Read(img gocv.Mat) bool {
	return true
}

type testWindow struct{}

func (w *testWindow) ShowImage(img gocv.Mat) { return }
