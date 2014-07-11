package opencv

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
)

type CameraDriver struct {
	gobot.Driver
	camera *cv.Capture
	Source interface{}
}

func NewCameraDriver(name string, source interface{}) *CameraDriver {
	c := &CameraDriver{
		Driver: *gobot.NewDriver(
			name,
			"CameraDriver",
		),
		Source: source,
	}

	c.AddEvent("frame")

	return c
}

func (c *CameraDriver) Start() bool {
	switch v := c.Source.(type) {
	case string:
		c.camera = cv.NewFileCapture(v)
	case int:
		c.camera = cv.NewCameraCapture(v)
	default:
		panic("unknown camera source")
	}

	go func() {
		for {
			if c.camera.GrabFrame() {
				image := c.camera.RetrieveFrame(1)
				if image != nil {
					gobot.Publish(c.Event("frame"), image)
				}
			}
		}
	}()
	return true
}

func (c *CameraDriver) Halt() bool { return true }
func (c *CameraDriver) Init() bool { return true }
