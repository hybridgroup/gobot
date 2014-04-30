package gpio

import (
	"github.com/hybridgroup/gobot"
)

type LedDriver struct {
	gobot.Driver
	Adaptor PwmDigitalWriter
	High    bool
}

func NewLedDriver(a PwmDigitalWriter, name, pin string) *LedDriver {
	return &LedDriver{
		Driver: gobot.Driver{
			Name: name,
			Pin:  pin,
			Commands: []string{
				"ToggleC",
				"OnC",
				"OffC",
				"BrightnessC",
			},
		},
		High:    false,
		Adaptor: a,
	}
}

func (l *LedDriver) Start() bool { return true }
func (l *LedDriver) Halt() bool  { return true }
func (l *LedDriver) Init() bool  { return true }

func (l *LedDriver) IsOn() bool {
	return l.High
}

func (l *LedDriver) IsOff() bool {
	return !l.IsOn()
}

func (l *LedDriver) On() bool {
	l.changeState(1)
	l.High = true
	return true
}

func (l *LedDriver) Off() bool {
	l.changeState(0)
	l.High = false
	return true
}

func (l *LedDriver) Toggle() {
	if l.IsOn() {
		l.Off()
	} else {
		l.On()
	}
}

func (l *LedDriver) Brightness(level byte) {
	l.Adaptor.PwmWrite(l.Pin, level)
}

func (l *LedDriver) changeState(level byte) {
	l.Adaptor.DigitalWrite(l.Pin, level)
}
