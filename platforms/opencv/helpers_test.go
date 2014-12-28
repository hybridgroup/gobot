package opencv

import cv "github.com/hybridgroup/go-opencv/opencv"

type testCapture struct{}

func (c *testCapture) RetrieveFrame(i int) *cv.IplImage {
	return &cv.IplImage{}
}

func (c *testCapture) GrabFrame() bool {
	return true
}

type testWindow struct{}

func (w *testWindow) ShowImage(i *cv.IplImage) { return }
