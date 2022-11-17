package system

import (
	"gobot.io/x/gobot"
)

// NewPWMPin returns a new system pwmPin.
func (a *Accesser) NewPWMPin(path string, pin int) gobot.PWMPinner {
	return newPWMPinSysfs(a.fs, path, pin)
}
