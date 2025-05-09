package adaptors

import (
	"fmt"

	"gobot.io/x/gobot/v2/system"
)

type PWMPinDefinition struct {
	Dir       string
	DirRegexp string
	Channel   int
}

type PWMPinDefinitions map[string]PWMPinDefinition

type PWMPinTranslator struct {
	sys            *system.Accesser
	pinDefinitions PWMPinDefinitions
}

// NewPWMPinTranslator creates a new instance of a PWM pin translator, suitable for the most cases.
func NewPWMPinTranslator(sys *system.Accesser, pinDefinitions PWMPinDefinitions) *PWMPinTranslator {
	return &PWMPinTranslator{sys: sys, pinDefinitions: pinDefinitions}
}

// Translate returns the sysfs path and channel for the given id.
func (pt *PWMPinTranslator) Translate(id string) (string, int, error) {
	pinInfo, ok := pt.pinDefinitions[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a PWM pin", id)
	}
	path, err := pinInfo.FindPWMDir(pt.sys)
	if err != nil {
		return "", -1, err
	}
	return path, pinInfo.Channel, nil
}

func (p PWMPinDefinition) FindPWMDir(sys *system.Accesser) (string, error) {
	items, _ := sys.Find(p.Dir, p.DirRegexp)
	if len(items) == 0 {
		return "", fmt.Errorf("No path found for PWM directory pattern, '%s' in path '%s'. See README.md for activation",
			p.DirRegexp, p.Dir)
	}

	dir := items[0]
	info, err := sys.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("Error (%v) on access '%s'", err, dir)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("The item '%s' is not a directory, which is not expected", dir)
	}

	return dir, nil
}
