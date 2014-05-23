package gpio

import (
	"github.com/hybridgroup/gobot"
)

type DirectPinDriver struct {
	gobot.Driver
	Adaptor DirectPin
}

func NewDirectPinDriver(a DirectPin, name string, pin string) *DirectPinDriver {
	return &DirectPinDriver{
		Driver: gobot.Driver{
			Name: name,
			Pin:  pin,
			Commands: []string{
				"DigitalReadC",
				"DigitalWriteC",
				"AnalogReadC",
				"AnalogWriteC",
				"PwmWriteC",
				"ServoWriteC",
			},
		},
		Adaptor: a,
	}
}

func (d *DirectPinDriver) Start() bool { return true }
func (d *DirectPinDriver) Halt() bool  { return true }
func (d *DirectPinDriver) Init() bool  { return true }

func (d *DirectPinDriver) DigitalRead() int {
	return d.Adaptor.DigitalRead(d.Pin)
}

func (d *DirectPinDriver) DigitalWrite(level byte) {
	d.Adaptor.DigitalWrite(d.Pin, level)
}

func (d *DirectPinDriver) AnalogRead() int {
	return d.Adaptor.AnalogRead(d.Pin)
}

func (d *DirectPinDriver) AnalogWrite(level byte) {
	d.Adaptor.AnalogWrite(d.Pin, level)
}

func (d *DirectPinDriver) PwmWrite(level byte) {
	d.Adaptor.PwmWrite(d.Pin, level)
}

func (d *DirectPinDriver) ServoWrite(level byte) {
	d.Adaptor.ServoWrite(d.Pin, level)
}
