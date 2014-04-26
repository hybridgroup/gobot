package gobotGPIO

import (
	"github.com/hybridgroup/gobot"
)

type DirectPinInterface interface {
	DigitalRead(string) int
	DigitalWrite(string, byte)
	AnalogRead(string) int
	AnalogWrite(string, byte)
	PwmWrite(string, byte)
	ServoWrite(string, byte)
}

type DirectPin struct {
	gobot.Driver
	Adaptor DirectPinInterface
}

func NewDirectPin(a DirectPinInterface) *DirectPin {
	b := new(DirectPin)
	b.Adaptor = a
	b.Events = make(map[string]chan interface{})
	b.Commands = []string{
		"DigitalReadC",
		"DigitalWriteC",
		"AnalogReadC",
		"AnalogWriteC",
		"PwmWriteC",
		"ServoWriteC",
	}
	return b
}

func (a *DirectPin) Start() bool { return true }
func (a *DirectPin) Halt() bool  { return true }
func (a *DirectPin) Init() bool  { return true }

func (a *DirectPin) DigitalRead() int {
	return a.Adaptor.DigitalRead(a.Pin)
}

func (a *DirectPin) DigitalWrite(level byte) {
	a.Adaptor.DigitalWrite(a.Pin, level)
}

func (a *DirectPin) AnalogRead() int {
	return a.Adaptor.AnalogRead(a.Pin)
}

func (a *DirectPin) AnalogWrite(level byte) {
	a.Adaptor.AnalogWrite(a.Pin, level)
}

func (a *DirectPin) PwmWrite(level byte) {
	a.Adaptor.PwmWrite(a.Pin, level)
}

func (a *DirectPin) ServoWrite(level byte) {
	a.Adaptor.ServoWrite(a.Pin, level)
}
