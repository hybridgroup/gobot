package opencv

import (
	"errors"

	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
)

var _ gobot.DriverInterface = (*CameraDriver)(nil)

type CameraDriver struct {
	gobot.Driver
	camera capture
	Source interface{}
	start  func(*CameraDriver) (err error)
}

// NewCameraDriver creates a new driver with specified name and source.
// It also creates a start function to either set camera as a File or Camera capture.
func NewCameraDriver(name string, source interface{}) *CameraDriver {
	c := &CameraDriver{
		Driver: *gobot.NewDriver(
			name,
			"CameraDriver",
		),
		Source: source,
		start: func(c *CameraDriver) (err error) {
			switch v := c.Source.(type) {
			case string:
				c.camera = cv.NewFileCapture(v)
			case int:
				c.camera = cv.NewCameraCapture(v)
			default:
				return errors.New("Unknown camera source")
			}
			return
		},
	}

	c.AddEvent("frame")

	return c
}

// Start initializes camera by grabbing a frame
// every `interval` and publishing an frame event
func (c *CameraDriver) Start() (err error) {
	if err = c.start(c); err != nil {
		return err
	}
	gobot.Every(c.Interval(), func() {
		if c.camera.GrabFrame() {
			image := c.camera.RetrieveFrame(1)
			if image != nil {
				gobot.Publish(c.Event("frame"), image)
			}
		}
	})
	return nil
}

// Halt stops camera driver
func (c *CameraDriver) Halt() error { return nil }
