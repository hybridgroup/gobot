package opencv

import cv "github.com/hybridgroup/go-opencv/opencv"

type capture interface {
	RetrieveFrame(int) *cv.IplImage
	GrabFrame() bool
}

type testCapture struct{}

func (c *testCapture) RetrieveFrame(i int) *cv.IplImage {
	return &cv.IplImage{}
}

func (c *testCapture) GrabFrame() bool {
	return true
}
