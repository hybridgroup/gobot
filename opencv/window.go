package gobotOpencv

import (
	"github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
)

type Window struct {
	gobot.Driver
	Adaptor *Opencv
	window  *opencv.Window
}

type WindowInterface interface {
}

func NewWindow(adaptor *Opencv) *Window {
	d := new(Window)
	d.Events = make(map[string]chan interface{})
	d.Adaptor = adaptor
	d.Commands = []string{}
	return d
}

func (me *Window) Start() bool {
	opencv.StartWindowThread()
	me.window = opencv.NewWindow(me.Name, opencv.CV_WINDOW_NORMAL)
	return true
}
func (me *Window) Halt() bool { return true }
func (me *Window) Init() bool { return true }

func (me *Window) ShowImage(image *opencv.IplImage) {
	me.window.ShowImage(image)
}
