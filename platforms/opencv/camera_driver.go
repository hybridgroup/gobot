package opencv

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
)

type CameraDriver struct {
	gobot.Driver
	camera capture
	Source interface{}
	start  func(*CameraDriver)
}

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

func (c *CameraDriver) Start() bool {
	c.start(c)
	gobot.Every(c.Interval(), func() {
		if c.camera.GrabFrame() {
			image := c.camera.RetrieveFrame(1)
			if image != nil {
				gobot.Publish(c.Event("frame"), image)
			}
		}
	})
	return true
}

func (c *CameraDriver) Halt() bool { return true }
