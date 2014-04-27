package opencv

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
)

type CameraDriver struct {
	gobot.Driver
	camera *cv.Capture
	Source string
}

func NewCameraDriver() *CameraDriver {
	return &CameraDriver{
		Driver: gobot.Driver{
			Commands: []string{},
			Events: map[string]chan interface{}{
				"Frame": make(chan interface{}, 0),
			},
		},
	}
}

func (c *CameraDriver) Start() bool {
	if c.Source != "" {
		c.camera = cv.NewFileCapture(c.Source)
	} else {
		c.camera = cv.NewCameraCapture(0)
	}
	go func() {
		for {
			if c.camera.GrabFrame() {
				image := c.camera.RetrieveFrame(1)
				if image != nil {
					gobot.Publish(c.Events["Frame"], image)
				}
			}
		}
	}()
	return true
}

func (c *CameraDriver) Halt() bool { return true }
func (c *CameraDriver) Init() bool { return true }
