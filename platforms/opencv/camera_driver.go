package opencv

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"time"
)

type CameraDriver struct {
	gobot.Driver
	camera capture
	Source interface{}
	start  func(*CameraDriver)
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
		start: func(c *CameraDriver) {
			switch v := c.Source.(type) {
			case string:
				c.camera = cv.NewFileCapture(v)
			case int:
				c.camera = cv.NewCameraCapture(v)
			default:
				panic("unknown camera source")
			}
		},
	}

	c.AddEvent("frame")

	return c
}

// Start initializes camera by grabbing a frame
// every `interval` and publishing an frame event
func (c *CameraDriver) Start() bool {
	c.start(c)
	go func() {
		for {
			if c.camera.GrabFrame() {
				image := c.camera.RetrieveFrame(1)
				if image != nil {
					gobot.Publish(c.Event("frame"), image)
				}
			}
			<-time.After(c.Interval())
		}
	}()
	return true
}

// Halt stops camera driver
func (c *CameraDriver) Halt() bool { return true }
