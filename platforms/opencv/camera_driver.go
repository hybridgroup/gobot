package opencv

import (
	"errors"

	"time"

	"github.com/hybridgroup/gobot"
	cv "github.com/lazywei/go-opencv/opencv"
)

type capture interface {
	RetrieveFrame(int) *cv.IplImage
	GrabFrame() bool
}

const (
	// Frame event
	Frame = "frame"
)

type CameraDriver struct {
	name     string
	camera   capture
	interval time.Duration
	Source   interface{}
	start    func(*CameraDriver) (err error)
	gobot.Eventer
}

// NewCameraDriver creates a new driver with specified name and source.
// It also creates a start function to either set camera as a File or Camera capture.
func NewCameraDriver(name string, source interface{}, v ...time.Duration) *CameraDriver {
	c := &CameraDriver{
		name:     name,
		Eventer:  gobot.NewEventer(),
		Source:   source,
		interval: 10 * time.Millisecond,
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

	if len(v) > 0 {
		c.interval = v[0]
	}

	c.AddEvent(Frame)

	return c
}

func (c *CameraDriver) Name() string                 { return c.name }
func (c *CameraDriver) Connection() gobot.Connection { return nil }

// Start initializes camera by grabbing a frame
// every `interval` and publishing an frame event
func (c *CameraDriver) Start() (errs []error) {
	if err := c.start(c); err != nil {
		return []error{err}
	}
	go func() {
		for {
			if c.camera.GrabFrame() {
				image := c.camera.RetrieveFrame(1)
				if image != nil {
					c.Publish(Frame, image)
				}
			}
			<-time.After(c.interval)
		}
	}()
	return
}

// Halt stops camera driver
func (c *CameraDriver) Halt() (errs []error) { return }
