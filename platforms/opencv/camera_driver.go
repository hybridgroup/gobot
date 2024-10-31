//go:build gocv
// +build gocv

package opencv

import (
	"errors"

	"gobot.io/x/gobot/v2"
	"gocv.io/x/gocv"
)

type capture interface {
	Read(img *gocv.Mat) bool
}

const (
	// Frame event
	Frame = "frame"
)

// CameraDriver is the Gobot Driver for the OpenCV camera
type CameraDriver struct {
	name   string
	camera capture
	Source interface{}
	start  func(*CameraDriver) error
	gobot.Eventer
}

// NewCameraDriver creates a new driver with specified source.
// It also creates a start function to either set camera as a File or Camera capture.
func NewCameraDriver(source interface{}) *CameraDriver {
	c := &CameraDriver{
		name:    "Camera",
		Eventer: gobot.NewEventer(),
		Source:  source,
		start: func(c *CameraDriver) error {
			switch v := c.Source.(type) {
			case string:
				c.camera, _ = gocv.VideoCaptureFile(v)
			case int:
				c.camera, _ = gocv.VideoCaptureDevice(v)
			default:
				return errors.New("Unknown camera source")
			}
			return nil
		},
	}

	c.AddEvent(Frame)

	return c
}

// Name returns the Driver name
func (c *CameraDriver) Name() string { return c.name }

// SetName sets the Driver name
func (c *CameraDriver) SetName(n string) { c.name = n }

// Connection returns the Driver's connection
func (c *CameraDriver) Connection() gobot.Connection { return nil }

// Start initializes camera by grabbing frames
func (c *CameraDriver) Start() error {
	if err := c.start(c); err != nil {
		return err
	}
	img := gocv.NewMat()
	go func() {
		for {
			if ok := c.camera.Read(&img); ok {
				c.Publish(Frame, img)
			}
		}
	}()
	return nil
}

// Halt stops camera driver
func (c *CameraDriver) Halt() error { return nil }
