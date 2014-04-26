package gobotGPIO

import (
	"github.com/hybridgroup/gobot"
)

type LedInterface interface {
	PwmWrite(string, byte)
	DigitalWrite(string, byte)
}

type Led struct {
	gobot.Driver
	Adaptor LedInterface
	High    bool
}

func NewLed(a LedInterface) *Led {
	l := new(Led)
	l.High = false
	l.Adaptor = a
	l.Commands = []string{
		"ToggleC",
		"OnC",
		"OffC",
		"BrightnessC",
	}
	return l
}

func (l *Led) Start() bool { return true }
func (l *Led) Halt() bool  { return true }
func (l *Led) Init() bool  { return true }

func (l *Led) IsOn() bool {
	return l.High
}

func (l *Led) IsOff() bool {
	return !l.IsOn()
}

func (l *Led) On() bool {
	l.changeState(1)
	l.High = true
	return true
}

func (l *Led) Off() bool {
	l.changeState(0)
	l.High = false
	return true
}

func (l *Led) Toggle() {
	if l.IsOn() {
		l.Off()
	} else {
		l.On()
	}
}

func (l *Led) Brightness(level byte) {
	l.Adaptor.PwmWrite(l.Pin, level)
}

func (l *Led) changeState(level byte) {
	l.Adaptor.DigitalWrite(l.Pin, level)
}
